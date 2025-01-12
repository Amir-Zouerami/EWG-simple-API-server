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
	ErrConflict            = errors.New("resource already exists")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		DeleteByID(ctx context.Context, postID int64) error
		UpdateByID(context.Context, *Post) error
		GetUserFeed(context.Context, int64, FeedPaginationQuery) ([]*FeedRecord, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		GetByPostID(ctx context.Context, postID int64) (*[]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followedID, userID int64) error
		Unfollow(ctx context.Context, unfollowedID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
