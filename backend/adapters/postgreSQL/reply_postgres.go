package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type ReplyPostgres struct {
	db *pgxpool.Pool
}

func NewReplyPostgres(db *pgxpool.Pool) repositories.ReplyRepository {
	return &ReplyPostgres{db: db}
}

func (r *ReplyPostgres) CreateReply(ctx context.Context, rep *entities.Reply) (int, error) {
	query := `INSERT INTO replies (thread_id, user_id, parent_id, body, is_deleted, created_at) VALUES ($1,$2,$3,$4,false,NOW()) RETURNING id`
	var id int
	err := r.db.QueryRow(ctx, query, rep.ThreadID, rep.UserID, rep.ParentID, rep.Body).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create reply: %w", err)
	}
	return id, nil
}

func (r *ReplyPostgres) GetRepliesByThread(ctx context.Context, threadID int) ([]entities.Reply, error) {
	query := `SELECT id, thread_id, user_id, parent_id, body, is_deleted, created_at, updated_at FROM replies WHERE thread_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, threadID)
	if err != nil {
		return nil, fmt.Errorf("get replies: %w", err)
	}
	defer rows.Close()
	var reps []entities.Reply
	for rows.Next() {
		var rep entities.Reply
		if err := rows.Scan(&rep.ID, &rep.ThreadID, &rep.UserID, &rep.ParentID, &rep.Body, &rep.IsDeleted, &rep.CreatedAt, &rep.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan reply: %w", err)
		}
		reps = append(reps, rep)
	}
	return reps, nil
}

func (r *ReplyPostgres) DeleteReply(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `UPDATE replies SET is_deleted = true, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete reply: %w", err)
	}
	return nil
}

func (r *ReplyPostgres) GetReplyByID(ctx context.Context, id int) (*entities.Reply, error) {
	var rep entities.Reply
	row := r.db.QueryRow(ctx, `SELECT id, thread_id, user_id, parent_id, body, is_deleted, created_at, updated_at FROM replies WHERE id = $1`, id)
	if err := row.Scan(&rep.ID, &rep.ThreadID, &rep.UserID, &rep.ParentID, &rep.Body, &rep.IsDeleted, &rep.CreatedAt, &rep.UpdatedAt); err != nil {
		return nil, fmt.Errorf("get reply by id: %w", err)
	}
	return &rep, nil
}
