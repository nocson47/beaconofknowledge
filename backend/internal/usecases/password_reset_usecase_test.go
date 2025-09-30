package usecases

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

// --- fakes ---
type fakeUserRepo struct {
	users map[int]*entities.User
}

func newFakeUserRepo() *fakeUserRepo                                             { return &fakeUserRepo{users: map[int]*entities.User{}} }
func (f *fakeUserRepo) GetAllUsers(ctx context.Context) ([]entities.User, error) { return nil, nil }
func (f *fakeUserRepo) CreateUser(ctx context.Context, user *entities.User) (int, error) {
	return 0, nil
}
func (f *fakeUserRepo) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}
func (f *fakeUserRepo) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	return nil, nil
}
func (f *fakeUserRepo) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}
func (f *fakeUserRepo) UpdateUser(ctx context.Context, user *entities.User) error {
	if _, ok := f.users[user.ID]; ok {
		f.users[user.ID] = user
	}
	return nil
}
func (f *fakeUserRepo) DeleteUser(ctx context.Context, id int) error { delete(f.users, id); return nil }

type fakePRRepo struct {
	byHash map[string]*entities.PasswordReset
}

func newFakePRRepo() *fakePRRepo { return &fakePRRepo{byHash: map[string]*entities.PasswordReset{}} }
func (f *fakePRRepo) Create(ctx context.Context, pr *entities.PasswordReset) (int, error) {
	f.byHash[pr.TokenHash] = pr
	pr.ID = 1
	return pr.ID, nil
}
func (f *fakePRRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordReset, error) {
	if pr, ok := f.byHash[tokenHash]; ok {
		return pr, nil
	}
	return nil, nil
}
func (f *fakePRRepo) MarkUsed(ctx context.Context, id int) error {
	for _, pr := range f.byHash {
		if pr.ID == id {
			pr.Used = true
		}
	}
	return nil
}
func (f *fakePRRepo) DeleteByUserID(ctx context.Context, userID int) error {
	for k, pr := range f.byHash {
		if pr.UserID == userID {
			delete(f.byHash, k)
		}
	}
	return nil
}

type fakeEmailSender struct{ lastURL string }

func (f *fakeEmailSender) SendResetEmail(ctx context.Context, toEmail string, resetURL string) error {
	f.lastURL = resetURL
	return nil
}

// --- tests ---
func TestPasswordResetUsecase_FullFlow(t *testing.T) {
	ctx := context.Background()
	urepo := newFakeUserRepo()
	urepo.users[1] = &entities.User{ID: 1, Email: "alice@example.com", Password: "oldhash"}
	prrepo := newFakePRRepo()
	email := &fakeEmailSender{}

	uc := NewPasswordResetUsecase(urepo, prrepo, email, 24*time.Hour)

	// Request reset
	if err := uc.RequestPasswordReset(ctx, "alice@example.com", "http://localhost:3000"); err != nil {
		t.Fatalf("RequestPasswordReset failed: %v", err)
	}
	if email.lastURL == "" {
		t.Fatalf("expected email URL to be set")
	}

	// extract token from URL
	u, err := url.Parse(email.lastURL)
	if err != nil {
		t.Fatalf("invalid reset URL: %v", err)
	}
	q := u.Query()
	token := q.Get("token")
	if token == "" {
		t.Fatalf("token missing in reset URL")
	}

	// Reset password
	if err := uc.ResetPassword(ctx, token, "newpass123"); err != nil {
		t.Fatalf("ResetPassword failed: %v", err)
	}

	// verify user's password was updated to a bcrypt hash that matches newpass123
	updated := urepo.users[1]
	if err := bcrypt.CompareHashAndPassword([]byte(updated.Password), []byte("newpass123")); err != nil {
		t.Fatalf("updated password does not match: %v", err)
	}

	// verify token was consumed: either marked used or deleted by DeleteByUserID
	usedFound := false
	for _, pr := range prrepo.byHash {
		if pr.Used {
			usedFound = true
		}
	}
	if !usedFound {
		// if not marked used, ensure token no longer exists (deleted)
		if len(prrepo.byHash) != 0 {
			t.Fatalf("expected token to be marked used or deleted, repository state: %#v", prrepo.byHash)
		}
	}
}

func TestPasswordResetUsecase_ExpiredToken(t *testing.T) {
	ctx := context.Background()
	urepo := newFakeUserRepo()
	urepo.users[1] = &entities.User{ID: 1, Email: "bob@example.com", Password: "old"}
	prrepo := newFakePRRepo()
	email := &fakeEmailSender{}

	// negative TTL to create already-expired token
	uc := NewPasswordResetUsecase(urepo, prrepo, email, -time.Hour)
	if err := uc.RequestPasswordReset(ctx, "bob@example.com", "http://localhost:3000"); err != nil {
		t.Fatalf("RequestPasswordReset failed: %v", err)
	}
	if email.lastURL == "" {
		t.Fatalf("expected email URL to be set")
	}
	u, _ := url.Parse(email.lastURL)
	token := u.Query().Get("token")

	if err := uc.ResetPassword(ctx, token, "whatever"); err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected expired token error, got: %v", err)
	}
}

func TestPasswordResetUsecase_UsedToken(t *testing.T) {
	ctx := context.Background()
	urepo := newFakeUserRepo()
	urepo.users[1] = &entities.User{ID: 1, Email: "carol@example.com", Password: "old"}
	prrepo := newFakePRRepo()
	email := &fakeEmailSender{}

	uc := NewPasswordResetUsecase(urepo, prrepo, email, 24*time.Hour)
	if err := uc.RequestPasswordReset(ctx, "carol@example.com", "http://localhost:3000"); err != nil {
		t.Fatalf("RequestPasswordReset failed: %v", err)
	}
	u, _ := url.Parse(email.lastURL)
	token := u.Query().Get("token")

	if err := uc.ResetPassword(ctx, token, "first"); err != nil {
		t.Fatalf("first ResetPassword failed: %v", err)
	}
	// second attempt should fail — implementation may either mark used or delete tokens; accept both "already used" or "invalid token"
	if err := uc.ResetPassword(ctx, token, "second"); err == nil {
		t.Fatalf("expected second ResetPassword to fail but succeeded")
	} else {
		if !strings.Contains(err.Error(), "already used") && !strings.Contains(err.Error(), "invalid token") {
			t.Fatalf("expected 'already used' or 'invalid token' error, got: %v", err)
		}
	}
}

func TestPasswordResetUsecase_InvalidToken(t *testing.T) {
	ctx := context.Background()
	urepo := newFakeUserRepo()
	urepo.users[1] = &entities.User{ID: 1, Email: "dan@example.com", Password: "old"}
	prrepo := newFakePRRepo()
	email := &fakeEmailSender{}

	uc := NewPasswordResetUsecase(urepo, prrepo, email, 24*time.Hour)
	// call ResetPassword with token that doesn't exist
	if err := uc.ResetPassword(ctx, "nonexistenttoken", "pw"); err == nil || !strings.Contains(err.Error(), "invalid token") {
		t.Fatalf("expected invalid token error, got: %v", err)
	}
}
