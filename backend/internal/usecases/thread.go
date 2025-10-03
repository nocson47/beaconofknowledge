// ...existing code...
package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

// ThreadService is the application-level port consumed by adapters (handlers).
type ThreadService interface {
	CreateThread(ctx context.Context, t *entities.Thread) (int, error)
	GetThreadByID(ctx context.Context, id int) (*entities.Thread, error)
	GetAllThreads(ctx context.Context) ([]*entities.Thread, error)
	UpdateThread(ctx context.Context, t *entities.Thread) error
	DeleteThread(ctx context.Context, id int) error
}

type threadService struct {
	repo repositories.ThreadRepository
}

func NewThreadService(repo repositories.ThreadRepository) ThreadService {
	return &threadService{repo: repo}
}

func (s *threadService) GetAllThreads(ctx context.Context) ([]*entities.Thread, error) {
	return s.repo.GetAllThreads(ctx)
}

func (s *threadService) CreateThread(ctx context.Context, t *entities.Thread) (int, error) {
	if t == nil {
		return 0, errors.New("thread is nil")
	}
	t.Title = strings.TrimSpace(t.Title)
	t.Body = strings.TrimSpace(t.Body)
	if t.Title == "" || t.Body == "" {
		return 0, errors.New("title and body are required")
	}
	id, err := s.repo.CreateThread(ctx, t)
	if err != nil {
		return 0, fmt.Errorf("create thread: %w", err)
	}
	return id, nil
}

func (s *threadService) GetThreadByID(ctx context.Context, id int) (*entities.Thread, error) {
	thread, err := s.repo.GetThreadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get thread by id %d: %w", id, err)
	}
	if thread == nil {
		return nil, errors.New("thread not found")
	}
	return thread, nil
}

func (s *threadService) UpdateThread(ctx context.Context, t *entities.Thread) error {
	if t == nil {
		return errors.New("thread is nil")
	}
	if strings.TrimSpace(t.Title) == "" || strings.TrimSpace(t.Body) == "" {
		return errors.New("title and body are required")
	}
	if err := s.repo.UpdateThread(ctx, t); err != nil {
		return fmt.Errorf("update thread: %w", err)
	}
	return nil
}

func (s *threadService) DeleteThread(ctx context.Context, id int) error {
	if err := s.repo.DeleteThread(ctx, id); err != nil {
		return fmt.Errorf("delete thread: %w", err)
	}
	return nil
}

// ...existing code...
