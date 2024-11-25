package socialmedia

import (
	"context"
	"time"
)

type PostLike struct {
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PostLikeService interface {
	CreatePostLike(ctx context.Context, postLike *PostLike) error
	FindPostLikes(ctx context.Context, filter PostLikeFilter) ([]*PostLike, int, error)
	DeletePostLike(ctx context.Context, id uint) error
}

type PostLikeFilter struct {
	PostID *uint `json:"post_id"`
	UserID *uint `json:"user_id"`

	Fill   bool
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
