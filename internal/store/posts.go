package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"comments"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (ps *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	  INSERT INTO posts (title, user_id, content, tags)
	  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := ps.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.UserID,
		post.Content,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ps *PostStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, version, created_at, updated_at, tags
		FROM posts WHERE id = $1 LIMIT 1;
	`

	var post Post

	err := ps.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (ps *PostStore) DeleteByID(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := ps.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (ps *PostStore) UpdateByID(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING version
	`
	err := ps.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
