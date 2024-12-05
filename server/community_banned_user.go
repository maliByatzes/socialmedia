package socialmedia

import (
	"context"
	"time"
)

type CBU struct {
	CommunityID uint      `json:"community_id"`
	UserID      uint      `json:"user_id"`
	BannedAt    time.Time `json:"banned_at"`
}

type CBUService interface {
	FindCBUs(ctx context.Context, filter CBUFilter) ([]*CBU, int, error)
	CreateCBU(ctx context.Context, cbu *CBU) error
	DeleteCBU(ctx context.Context, id uint) error
}

type CBUFilter struct {
	CommunityID *uint `json:"community_id"`
	UserID      *uint `json:"user_id"`
}
