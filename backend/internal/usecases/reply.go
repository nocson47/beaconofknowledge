package usecases

import (
	"context"
	"fmt"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type ReplyService interface {
	CreateReply(ctx context.Context, r *entities.Reply) (int, error)
	GetRepliesByThread(ctx context.Context, threadID int) ([]entities.Reply, error)
	GetReplyByID(ctx context.Context, id int) (*entities.Reply, error)
	// DeleteReply enforces authorization: admins can delete any reply, users can delete their own
	DeleteReply(ctx context.Context, id int, actorUserID int, isAdmin bool) error
}

type replyService struct {
	repo repositories.ReplyRepository
}

func NewReplyService(repo repositories.ReplyRepository) ReplyService {
	return &replyService{repo: repo}
}

func (s *replyService) CreateReply(ctx context.Context, r *entities.Reply) (int, error) {
	if r == nil {
		return 0, fmt.Errorf("reply is nil")
	}
	if r.Body == "" {
		return 0, fmt.Errorf("body is empty")
	}
	return s.repo.CreateReply(ctx, r)
}

func (s *replyService) GetRepliesByThread(ctx context.Context, threadID int) ([]entities.Reply, error) {
	return s.repo.GetRepliesByThread(ctx, threadID)
}

func (s *replyService) GetReplyByID(ctx context.Context, id int) (*entities.Reply, error) {
	return s.repo.GetReplyByID(ctx, id)
}

func (s *replyService) DeleteReply(ctx context.Context, id int, actorUserID int, isAdmin bool) error {
	rep, err := s.repo.GetReplyByID(ctx, id)
	if err != nil {
		return err
	}
	if rep == nil {
		return fmt.Errorf("reply not found")
	}
	// If actor is admin, allow delete. Otherwise, only allow if actorUserID == rep.UserID
	if !isAdmin && actorUserID != rep.UserID {
		return fmt.Errorf("forbidden: cannot delete others' replies")
	}
	return s.repo.DeleteReply(ctx, id)
}
