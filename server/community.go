package socialmedia

import (
	"context"
	"time"
)

type Community struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Banner      string    `json:"banner"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Community) Validate() error {
	if c.Name == "" {
		return Errorf(EINVALID, "Name is required.")
	}

	return nil
}

type CommunityService interface {
	FindCommunityByID(ctx context.Context, id uint) (*Community, error)
	FindCommunities(ctx context.Context, filter CommunityFilter) ([]*Community, int, error)
	CreateCommunity(ctx context.Context, com *Community) error
	UpdateCommunity(ctx context.Context, id uint, upd CommunityUpdate) (*Community, error)
	DeleteCommunity(ctx context.Context, id uint) error
}

type CommunityFilter struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CommunityUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Banner      *string `json:"banner"`
}
