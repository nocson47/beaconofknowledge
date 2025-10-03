package usecases

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

// EmailSender is an interface for sending emails; can be implemented by console or SMTP adapters.
type EmailSender interface {
	SendResetEmail(ctx context.Context, toEmail string, resetURL string) error
}

// PasswordResetUsecase contains dependencies for password reset flow
type PasswordResetUsecase struct {
	prRepo   repositories.PasswordResetRepository
	users    repositories.UserRepository
	email    EmailSender
	tokenTTL time.Duration
}

// NewPasswordResetUsecase constructs the usecase
func NewPasswordResetUsecase(users repositories.UserRepository, prRepo repositories.PasswordResetRepository, email EmailSender, ttl time.Duration) *PasswordResetUsecase {
	return &PasswordResetUsecase{prRepo: prRepo, users: users, email: email, tokenTTL: ttl}
}

// RequestPasswordReset handles generating a token and storing its hash; sends reset link via email sender
func (uc *PasswordResetUsecase) RequestPasswordReset(ctx context.Context, email string, baseURL string) error {
	user, err := uc.users.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %w", err)
	}
	if user == nil {
		// avoid leaking whether email exists: return nil (no-op)
		return nil
	}

	// create secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// hash token before storing
	h := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(h[:])

	now := time.Now().UTC()
	pr := &entities.PasswordReset{
		UserID:    user.ID,
		TokenHash: tokenHash,
		CreatedAt: now,
		ExpiresAt: now.Add(uc.tokenTTL),
		Used:      false,
	}

	if _, err := uc.prRepo.Create(ctx, pr); err != nil {
		return fmt.Errorf("failed to store password reset: %w", err)
	}

	// build reset URL (baseURL should include scheme and host, e.g. https://example.com)
	resetURL := fmt.Sprintf("%s/reset?token=%s", baseURL, token)

	// send email (or console)
	if err := uc.email.SendResetEmail(ctx, user.Email, resetURL); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword verifies token, updates user's password and marks token used
func (uc *PasswordResetUsecase) ResetPassword(ctx context.Context, token string, newPassword string) error {
	// hash incoming token
	h := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(h[:])

	pr, err := uc.prRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find token: %w", err)
	}
	if pr == nil {
		return fmt.Errorf("invalid token")
	}
	if pr.Used {
		return fmt.Errorf("token already used")
	}
	if time.Now().UTC().After(pr.ExpiresAt) {
		return fmt.Errorf("token expired")
	}

	user, err := uc.users.GetUserByID(ctx, pr.UserID)
	if err != nil {
		return fmt.Errorf("failed to load user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// bcrypt hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashed)
	if err := uc.users.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	// mark token used
	if err := uc.prRepo.MarkUsed(ctx, pr.ID); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	// optional: delete other tokens for this user
	_ = uc.prRepo.DeleteByUserID(ctx, user.ID)

	return nil
}
