package socialmedia

import (
	"context"
	"time"
)

type Post struct {
	ID          uint      `json:"id"`
	Content     string    `json:"content"`
	FileURL     string    `json:"file_url"`
	CommunityID uint      `json:"community_id"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostService interface {
	FindPostByID(ctx context.Context, id uint) (*Post, error)
	FindPosts(ctx context.Context, filter PostFilter) ([]*Post, int, error)
	CreatePost(ctx context.Context, post *Post) error
	UpdatePost(ctx context.Context, id uint, upd PostUpdate) (*Post, error)
	DeletePost(ctx context.Context, id uint) error
}

type PostFilter struct {
	ID          *uint `json:"id"`
	CommunityID *uint `json:"community_id"`
	UserID      *uint `json:"user_id"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PostUpdate struct {
	Content *string `json:"content"`
	FileURL *string `json:"file_url"`
}
