package socialmedia

import (
	"context"
	"time"
)

type CommunityMember struct {
	CommunityID uint      `json:"community_id"`
	UserID      uint      `json:"user_id"`
	IsModerator bool      `json:"is_moderator"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommunityMemberService interface {
	FindCommunityMembers(ctx context.Context, filter CommunityMemberFilter) ([]*CommunityMember, int, error)
	CreateCommunityMember(ctx context.Context, cm *CommunityMember) error
	UpdateCommunityMember(ctx context.Context, id uint, upd CommunityMemberUpdate) (*CommunityMember, error)
	DeleteCommunityMember(ctx context.Context, id uint) error
}

type CommunityMemberFilter struct {
	CommunityID *uint `json:"community_id"`
	UserID      *uint `json:"user_id"`
	IsModerator *bool `json:"is_moderator"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CommunityMemberUpdate struct {
	IsModerator *bool `json:"is_moderator"`
}
