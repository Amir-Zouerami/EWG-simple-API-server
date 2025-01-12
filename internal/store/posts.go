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
	User      User      `json:"user"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"comments"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type FeedRecord struct {
	Post
	CommentsCount int `json:"comments_count"`
}

func (ps *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	  INSERT INTO posts (title, user_id, content, tags)
	  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

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

func (ps *PostStore) GetUserFeed(ctx context.Context, userID int64, fq FeedPaginationQuery) ([]*FeedRecord, error) {
	query := `
	SELECT
	p.id,
	p.user_id,
	p.title,
	p.content,
	p.created_at,
	p.version,
	p.tags,
	u.username,
	COALESCE(comment_counts.comment_count, 0) AS comment_count
	FROM
	posts p
	LEFT JOIN (
		SELECT post_id, COUNT(*) AS comment_count
		FROM comments
		GROUP BY post_id
	) comment_counts ON p.id = comment_counts.post_id
	LEFT JOIN users u ON p.user_id = u.id
	WHERE
	p.user_id = $1
	OR p.user_id IN (
		SELECT f.follower_id
		FROM followers f
		WHERE f.user_id = $1
	)
	ORDER BY p.created_at ` + fq.Sort + ` LIMIT $2 OFFSET $3
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	rows, err := ps.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feedRecords []*FeedRecord

	for rows.Next() {
		var record FeedRecord

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Title,
			&record.Content,
			&record.CreatedAt,
			&record.Version,
			pq.Array(&record.Tags),
			&record.User.Username,
			&record.CommentsCount,
		)

		if err != nil {
			return nil, err
		}

		feedRecords = append(feedRecords, &record)
	}

	return feedRecords, nil
}
