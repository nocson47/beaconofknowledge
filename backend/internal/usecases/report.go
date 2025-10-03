package usecases

import (
	"context"
	"fmt"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type ReportService interface {
	CreateReport(ctx context.Context, r *entities.Report) (string, error)
	GetReports(ctx context.Context, kind *string) ([]*entities.Report, error)
	UpdateReportStatus(ctx context.Context, id string, status string, resolvedBy *int) error
}

type reportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) CreateReport(ctx context.Context, r *entities.Report) (string, error) {
	if r == nil {
		return "", fmt.Errorf("report is nil")
	}
	if r.Kind != "thread" && r.Kind != "user" {
		return "", fmt.Errorf("invalid kind")
	}
	if r.TargetID == 0 {
		return "", fmt.Errorf("target_id required")
	}
	return s.repo.CreateReport(ctx, r)
}

func (s *reportService) GetReports(ctx context.Context, kind *string) ([]*entities.Report, error) {
	return s.repo.GetReports(ctx, kind)
}

func (s *reportService) UpdateReportStatus(ctx context.Context, id string, status string, resolvedBy *int) error {
	if status != "open" && status != "resolved" && status != "dismissed" {
		return fmt.Errorf("invalid status")
	}
	return s.repo.UpdateReportStatus(ctx, id, status, resolvedBy)
}
