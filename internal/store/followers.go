package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (store *FollowerStore) Follow(ctx context.Context, followedID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	query := `
	  INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)
	`

	_, err := store.db.ExecContext(ctx, query, userID, followedID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return err
}

func (store *FollowerStore) Unfollow(ctx context.Context, unfollowedID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	query := `
	  DELETE FROM followers WHERE user_id = $1 AND follower_id = $2
	`

	_, err := store.db.ExecContext(ctx, query, userID, unfollowedID)
	return err
}
