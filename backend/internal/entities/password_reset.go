package entities

import "time"

// PasswordReset represents a password reset token record
type PasswordReset struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	TokenHash string    `db:"token_hash" json:"token_hash"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Used      bool      `db:"used" json:"used"`
}
