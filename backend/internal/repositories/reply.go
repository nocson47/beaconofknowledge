package repositories

import (
	"context"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

type ReplyRepository interface {
	CreateReply(ctx context.Context, r *entities.Reply) (int, error)
	GetRepliesByThread(ctx context.Context, threadID int) ([]entities.Reply, error)
	GetReplyByID(ctx context.Context, id int) (*entities.Reply, error)
	DeleteReply(ctx context.Context, id int) error
}
