package repositories

import (
	"context"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) (int, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id int) error
}
