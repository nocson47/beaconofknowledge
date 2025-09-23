package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type ThreadPostgres struct {
	db *pgxpool.Pool
}

func NewThreadPostgres(db *pgxpool.Pool) repositories.ThreadRepository {
	return &ThreadPostgres{db: db}
}

// CreateThread inserts thread and links tags (creates tags if needed).
func (r *ThreadPostgres) CreateThread(ctx context.Context, thread *entities.Thread) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	insertThread := `INSERT INTO threads (user_id, title, body, is_locked, is_deleted, created_at, updated_at)
                     VALUES ($1,$2,$3,$4,$5,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) RETURNING id`
	var id int
	err = tx.QueryRow(ctx, insertThread, thread.UserID, thread.Title, thread.Body, thread.IsLocked, thread.IsDeleted).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert thread: %w", err)
	}

	// handle tags (optional)
	for _, tag := range thread.Tags {
		var tagID int
		// upsert tag
		err = tx.QueryRow(ctx, `INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, tag).Scan(&tagID)
		if err != nil {
			// if RETURNING fails due to conflict handling, fetch id
			err = tx.QueryRow(ctx, `SELECT id FROM tags WHERE name = $1`, tag).Scan(&tagID)
			if err != nil {
				return 0, fmt.Errorf("ensure tag: %w", err)
			}
		}
		_, err = tx.Exec(ctx, `INSERT INTO thread_tags (thread_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, tagID)
		if err != nil {
			return 0, fmt.Errorf("link tag: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return id, nil
}

func (r *ThreadPostgres) GetThreadByID(ctx context.Context, id int) (*entities.Thread, error) {
	query := `
    SELECT t.id, t.user_id, t.title, t.body, t.is_locked, t.is_deleted, t.upvotes, t.downvotes, t.created_at, t.updated_at,
	    COALESCE(array_agg(tags.name) FILTER (WHERE tags.name IS NOT NULL), '{}') AS tags
    FROM threads t
    LEFT JOIN thread_tags tt ON tt.thread_id = t.id
    LEFT JOIN tags ON tags.id = tt.tag_id
    WHERE t.id = $1
    GROUP BY t.id;
    `
	var thread entities.Thread
	var tags []string
	if err := r.db.QueryRow(ctx, query, id).Scan(&thread.ID, &thread.UserID, &thread.Title, &thread.Body, &thread.IsLocked, &thread.IsDeleted, &thread.Upvotes, &thread.Downvotes, &thread.CreatedAt, &thread.UpdatedAt, &tags); err != nil {
		// use pgxpool/errors handling at caller
		return nil, fmt.Errorf("failed to retrieve thread: %w", err)
	}
	thread.Tags = tags
	return &thread, nil
}

func (r *ThreadPostgres) UpdateThread(ctx context.Context, thread *entities.Thread) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	update := `UPDATE threads SET title=$1, body=$2, is_locked=$3, is_deleted=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5`
	_, err = tx.Exec(ctx, update, thread.Title, thread.Body, thread.IsLocked, thread.IsDeleted, thread.ID)
	if err != nil {
		return fmt.Errorf("update thread: %w", err)
	}

	// replace tags: simple approach - delete existing links, then ensure tags and link them
	_, err = tx.Exec(ctx, `DELETE FROM thread_tags WHERE thread_id = $1`, thread.ID)
	if err != nil {
		return fmt.Errorf("clear thread tags: %w", err)
	}
	for _, tag := range thread.Tags {
		var tagID int
		err = tx.QueryRow(ctx, `INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, tag).Scan(&tagID)
		if err != nil {
			err = tx.QueryRow(ctx, `SELECT id FROM tags WHERE name = $1`, tag).Scan(&tagID)
			if err != nil {
				return fmt.Errorf("ensure tag: %w", err)
			}
		}
		_, err = tx.Exec(ctx, `INSERT INTO thread_tags (thread_id, tag_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, thread.ID, tagID)
		if err != nil {
			return fmt.Errorf("link tag: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *ThreadPostgres) DeleteThread(ctx context.Context, id int) error {
	query := `UPDATE threads SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}
	return nil
}

func (r *ThreadPostgres) GetAllThreads(ctx context.Context) ([]*entities.Thread, error) {
	query := `
	SELECT t.id, t.user_id, t.title, t.body, t.is_locked, t.is_deleted, t.upvotes, t.downvotes, t.created_at, t.updated_at,
		COALESCE(array_agg(tags.name) FILTER (WHERE tags.name IS NOT NULL), '{}') AS tags
	FROM threads t
	LEFT JOIN thread_tags tt ON tt.thread_id = t.id
	LEFT JOIN tags ON tags.id = tt.tag_id
	GROUP BY t.id;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve threads: %w", err)
	}
	defer rows.Close()

	var threads []*entities.Thread
	for rows.Next() {
		var th entities.Thread
		var tags []string
		if err := rows.Scan(&th.ID, &th.UserID, &th.Title, &th.Body, &th.IsLocked, &th.IsDeleted, &th.Upvotes, &th.Downvotes, &th.CreatedAt, &th.UpdatedAt, &tags); err != nil {
			return nil, fmt.Errorf("failed to scan thread: %w", err)
		}
		th.Tags = tags
		threads = append(threads, &th)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return threads, nil
}
