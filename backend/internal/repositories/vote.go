package repositories

import (
	"context"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

type VoteRepository interface {
	CreateVote(ctx context.Context, v *entities.Vote) (int, error)
	GetVoteByID(ctx context.Context, id int) (*entities.Vote, error)
	GetVotesForThread(ctx context.Context, threadID int) ([]entities.Vote, error)
	GetVotesForReply(ctx context.Context, replyID int) ([]entities.Vote, error)
	DeleteVote(ctx context.Context, id int) error
	GetVoteCountsForThread(ctx context.Context, threadID int) (up int, down int, err error)
}
