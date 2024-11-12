package socialmedia

import (
	"context"
	"time"
)

type Token struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TokenService interface {
	FindTokenByID(ctx context.Context, id uint) (*Token, error)
	FindTokens(ctx context.Context, filter TokenFilter) ([]*Token, int, error)
	CreateToken(ctx context.Context, token *Token) error
}

type TokenFilter struct {
	ID     *uint `json:"id"`
	UserID *uint `json:"user_id"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
