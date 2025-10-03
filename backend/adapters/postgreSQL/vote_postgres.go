package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
)

type VotePostgres struct {
	db *pgxpool.Pool
}

func NewVotePostgres(db *pgxpool.Pool) repositories.VoteRepository {
	return &VotePostgres{db: db}
}

func (r *VotePostgres) CreateVote(ctx context.Context, v *entities.Vote) (int, error) {
	// Transactional insert/update: keep votes and thread counters consistent
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		// if still in progress, rollback
		_ = tx.Rollback(ctx)
	}()

	var existingValue int
	var voteID int
	// try to select existing vote for this user/thread
	if v.ThreadID != nil {
		err = tx.QueryRow(ctx, `SELECT id, value FROM votes WHERE user_id=$1 AND thread_id=$2 FOR UPDATE`, v.UserID, *v.ThreadID).Scan(&voteID, &existingValue)
		if err != nil && err.Error() == "no rows in result set" {
			// no existing vote
			voteID = 0
			err = nil
		}
	} else {
		// replies not yet updating thread counters here; try to select by reply_id
		if v.ReplyID != nil {
			err = tx.QueryRow(ctx, `SELECT id, value FROM votes WHERE user_id=$1 AND reply_id=$2 FOR UPDATE`, v.UserID, *v.ReplyID).Scan(&voteID, &existingValue)
			if err != nil && err.Error() == "no rows in result set" {
				voteID = 0
				err = nil
			}
		}
	}
	if err != nil {
		return 0, fmt.Errorf("select existing vote: %w", err)
	}

	deltaUp, deltaDown := 0, 0

	if voteID == 0 {
		// insert
		insertQ := `INSERT INTO votes (user_id, thread_id, reply_id, value, created_at) VALUES ($1,$2,$3,$4,NOW()) RETURNING id`
		err = tx.QueryRow(ctx, insertQ, v.UserID, v.ThreadID, v.ReplyID, v.Value).Scan(&voteID)
		if err != nil {
			return 0, fmt.Errorf("insert vote: %w", err)
		}
		if v.Value == 1 {
			deltaUp = 1
		} else if v.Value == -1 {
			deltaDown = 1
		}
	} else {
		// existing vote
		if existingValue == v.Value {
			// no-op
		} else {
			// update
			_, err = tx.Exec(ctx, `UPDATE votes SET value=$1 WHERE id=$2`, v.Value, voteID)
			if err != nil {
				return 0, fmt.Errorf("update vote: %w", err)
			}
			// compute deltas
			if existingValue == 1 && v.Value == -1 {
				deltaUp = -1
				deltaDown = 1
			} else if existingValue == -1 && v.Value == 1 {
				deltaUp = 1
				deltaDown = -1
			}
		}
	}

	// update thread counters if this vote is for a thread
	if v.ThreadID != nil && (deltaUp != 0 || deltaDown != 0) {
		_, err = tx.Exec(ctx, `UPDATE threads SET upvotes = upvotes + $1, downvotes = downvotes + $2 WHERE id = $3`, deltaUp, deltaDown, *v.ThreadID)
		if err != nil {
			return 0, fmt.Errorf("update thread counters: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return voteID, nil
}

func (r *VotePostgres) GetVoteByID(ctx context.Context, id int) (*entities.Vote, error) {
	query := `SELECT id, user_id, thread_id, reply_id, value, created_at FROM votes WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var v entities.Vote
	if err := row.Scan(&v.ID, &v.UserID, &v.ThreadID, &v.ReplyID, &v.Value, &v.CreatedAt); err != nil {
		return nil, fmt.Errorf("get vote: %w", err)
	}
	return &v, nil
}

func (r *VotePostgres) GetVotesForThread(ctx context.Context, threadID int) ([]entities.Vote, error) {
	query := `SELECT id, user_id, thread_id, reply_id, value, created_at FROM votes WHERE thread_id = $1`
	rows, err := r.db.Query(ctx, query, threadID)
	if err != nil {
		return nil, fmt.Errorf("get votes for thread: %w", err)
	}
	defer rows.Close()
	var votes []entities.Vote
	for rows.Next() {
		var v entities.Vote
		if err := rows.Scan(&v.ID, &v.UserID, &v.ThreadID, &v.ReplyID, &v.Value, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan vote: %w", err)
		}
		votes = append(votes, v)
	}
	return votes, nil
}

func (r *VotePostgres) GetVotesForReply(ctx context.Context, replyID int) ([]entities.Vote, error) {
	query := `SELECT id, user_id, thread_id, reply_id, value, created_at FROM votes WHERE reply_id = $1`
	rows, err := r.db.Query(ctx, query, replyID)
	if err != nil {
		return nil, fmt.Errorf("get votes for reply: %w", err)
	}
	defer rows.Close()
	var votes []entities.Vote
	for rows.Next() {
		var v entities.Vote
		if err := rows.Scan(&v.ID, &v.UserID, &v.ThreadID, &v.ReplyID, &v.Value, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan vote: %w", err)
		}
		votes = append(votes, v)
	}
	return votes, nil
}

func (r *VotePostgres) DeleteVote(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM votes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete vote: %w", err)
	}
	return nil
}

func (r *VotePostgres) GetVoteCountsForThread(ctx context.Context, threadID int) (int, int, error) {
	query := `SELECT SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END) AS upvotes, SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END) AS downvotes FROM votes WHERE thread_id = $1`
	var up, down int
	if err := r.db.QueryRow(ctx, query, threadID).Scan(&up, &down); err != nil {
		return 0, 0, fmt.Errorf("get vote counts: %w", err)
	}
	return up, down, nil
}
