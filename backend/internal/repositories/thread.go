package repositories

import (
	"context"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

type ThreadRepository interface {
	CreateThread(ctx context.Context, thread *entities.Thread) (int, error)
	GetThreadByID(ctx context.Context, id int) (*entities.Thread, error)
	GetAllThreads(ctx context.Context) ([]*entities.Thread, error)
	UpdateThread(ctx context.Context, thread *entities.Thread) error
	DeleteThread(ctx context.Context, id int) error
}
