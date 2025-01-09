package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QUERY_TIMEOUT_DURATION = time.Second * 5
	ErrNotFound            = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		DeleteByID(ctx context.Context, postID int64) error
		UpdateByID(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		GetByPostID(ctx context.Context, postID int64) (*[]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
