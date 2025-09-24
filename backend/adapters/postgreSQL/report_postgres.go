package postgressql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type ReportPostgres struct {
	db *pgxpool.Pool
}

func NewReportPostgres(db *pgxpool.Pool) repositories.ReportRepository {
	return &ReportPostgres{db: db}
}

func (r *ReportPostgres) CreateReport(ctx context.Context, rep *entities.Report) (string, error) {
	query := `INSERT INTO reports (reporter_id, kind, target_id, reason, status, created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	var id int
	err := r.db.QueryRow(ctx, query, rep.ReporterID, rep.Kind, rep.TargetID, rep.Reason, rep.Status, time.Now()).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("create report: %w", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func (r *ReportPostgres) GetReports(ctx context.Context, kind *string) ([]*entities.Report, error) {
	var rows pgx.Rows
	var err error
	if kind != nil {
		rows, err = r.db.Query(ctx, `SELECT id, reporter_id, kind, target_id, reason, status, created_at, resolved_by, resolved_at FROM reports WHERE kind=$1 ORDER BY created_at DESC`, *kind)
	} else {
		rows, err = r.db.Query(ctx, `SELECT id, reporter_id, kind, target_id, reason, status, created_at, resolved_by, resolved_at FROM reports ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("query reports: %w", err)
	}
	defer rows.Close()

	var out []*entities.Report
	for rows.Next() {
		var rep entities.Report
		var resolvedAt *time.Time
		var reporterID *int
		var resolvedBy *int
		if err := rows.Scan(&rep.ID, &reporterID, &rep.Kind, &rep.TargetID, &rep.Reason, &rep.Status, &rep.CreatedAt, &resolvedBy, &resolvedAt); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		rep.ReporterID = reporterID
		rep.ResolvedBy = resolvedBy
		rep.ResolvedAt = resolvedAt
		out = append(out, &rep)
	}
	return out, nil
}

func (r *ReportPostgres) UpdateReportStatus(ctx context.Context, id string, status string, resolvedBy *int) error {
	// convert id string to int
	iid, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	var resolvedAt interface{}
	if resolvedBy != nil {
		resolvedAt = time.Now()
	} else {
		resolvedAt = nil
	}
	_, err = r.db.Exec(ctx, `UPDATE reports SET status=$1, resolved_by=$2, resolved_at=$3 WHERE id=$4`, status, resolvedBy, resolvedAt, iid)
	if err != nil {
		return fmt.Errorf("update report: %w", err)
	}
	return nil
}
