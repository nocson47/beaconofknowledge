package repositories

import (
	"context"

	"github.com/nocson47/beaconofknowledge/internal/entities"
)

type ReportRepository interface {
	CreateReport(ctx context.Context, r *entities.Report) (string, error)
	GetReports(ctx context.Context, kind *string) ([]*entities.Report, error)
	UpdateReportStatus(ctx context.Context, id string, status string, resolvedBy *int) error
}
