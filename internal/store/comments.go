package store

import (
	"context"
	"database/sql"
)

type CommentStore struct {
	db *sql.DB
}

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	User      User   `json:"user"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) (*[]Comment, error) {
	query := `
	  SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id FROM comments c
	  JOIN users u on u.id = c.user_id
	  WHERE c.post_id = $1
	  ORDER BY c.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment

		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.User.Username, &c.User.ID, &c.CreatedAt)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return &comments, nil
}
