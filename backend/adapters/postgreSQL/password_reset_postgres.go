package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type PasswordResetPostgres struct {
	db *pgxpool.Pool
}

func NewPasswordResetPostgres(db *pgxpool.Pool) repositories.PasswordResetRepository {
	return &PasswordResetPostgres{db: db}
}

func (p *PasswordResetPostgres) Create(ctx context.Context, pr *entities.PasswordReset) (int, error) {
	query := `INSERT INTO password_resets (user_id, token_hash, created_at, expires_at, used) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	var id int
	err := p.db.QueryRow(ctx, query, pr.UserID, pr.TokenHash, pr.CreatedAt, pr.ExpiresAt, pr.Used).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create password_reset: %w", err)
	}
	return id, nil
}

func (p *PasswordResetPostgres) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordReset, error) {
	query := `SELECT id, user_id, token_hash, created_at, expires_at, used FROM password_resets WHERE token_hash = $1`
	row := p.db.QueryRow(ctx, query, tokenHash)
	var pr entities.PasswordReset
	if err := row.Scan(&pr.ID, &pr.UserID, &pr.TokenHash, &pr.CreatedAt, &pr.ExpiresAt, &pr.Used); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find password_reset: %w", err)
	}
	return &pr, nil
}

func (p *PasswordResetPostgres) MarkUsed(ctx context.Context, id int) error {
	query := `UPDATE password_resets SET used = true WHERE id = $1`
	_, err := p.db.Exec(ctx, query, id)
	return err
}

func (p *PasswordResetPostgres) DeleteByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM password_resets WHERE user_id = $1`
	_, err := p.db.Exec(ctx, query, userID)
	return err
}
