package socialmedia

import (
	"context"
	"time"
)

type Relationship struct {
	ID          uint      `json:"id"`
	FollowerID  uint      `json:"follower_id"`
	FollowingID uint      `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RelationshipService interface {
	FindRelationshipByID(ctx context.Context, id uint) (*Relationship, error)
	FindRelationships(ctx context.Context, filter RelationshipFilter) ([]*Relationship, int, error)
	CreateRelationship(ctx context.Context, r *Relationship) error
	DeleteRelationship(ctx context.Context, id uint) error
}

type RelationshipFilter struct {
	ID          *uint `json:"id"`
	FollowerID  *uint `json:"follower_id"`
	FollowingID *uint `json:"following_id"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
