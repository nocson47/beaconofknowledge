package repositories

import (
	"context"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, pr *entities.PasswordReset) (int, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordReset, error)
	MarkUsed(ctx context.Context, id int) error
	DeleteByUserID(ctx context.Context, userID int) error
}
