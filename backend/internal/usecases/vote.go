package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type VoteService interface {
	CreateVote(ctx context.Context, v *entities.Vote) (int, error)
	GetVoteByID(ctx context.Context, id int) (*entities.Vote, error)
	GetVoteCountsForThread(ctx context.Context, threadID int) (int, int, error)
}

type voteService struct {
	repo repositories.VoteRepository
}

func NewVoteService(repo repositories.VoteRepository) VoteService {
	return &voteService{repo: repo}
}

func (s *voteService) CreateVote(ctx context.Context, v *entities.Vote) (int, error) {
	if v == nil {
		return 0, errors.New("vote is nil")
	}
	if v.Value != 1 && v.Value != -1 {
		return 0, errors.New("invalid vote value")
	}
	// delegate to repository (repo can enforce unique constraint)
	id, err := s.repo.CreateVote(ctx, v)
	if err != nil {
		return 0, fmt.Errorf("create vote: %w", err)
	}
	return id, nil
}

func (s *voteService) GetVoteByID(ctx context.Context, id int) (*entities.Vote, error) {
	return s.repo.GetVoteByID(ctx, id)
}

func (s *voteService) GetVoteCountsForThread(ctx context.Context, threadID int) (int, int, error) {
	return s.repo.GetVoteCountsForThread(ctx, threadID)
}
