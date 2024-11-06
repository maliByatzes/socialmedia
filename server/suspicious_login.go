package socialmedia

import (
	"context"
	"time"
)

type SuspiciousLogin struct {
	ID                 uint      `json:"id"`
	UserID             uint      `json:"user_id"`
	Email              string    `json:"email"`
	IP                 string    `json:"ip"`
	Country            string    `json:"country"`
	City               string    `json:"city"`
	Browser            string    `json:"browser"`
	Platform           string    `json:"platform"`
	OS                 string    `json:"os"`
	Device             string    `json:"device"`
	DeviceType         string    `json:"device_type"`
	UnverifiedAttempts int       `json:"unverified_attempts"`
	IsTrusted          bool      `json:"is_trusted"`
	IsBlocked          bool      `json:"is_blocked"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type SuspiciousLoginService interface {
	FindSLByID(ctx context.Context, id uint) (*SuspiciousLogin, error)
	FindSLByUserID(ctx context.Context, userID uint) (*SuspiciousLogin, error)
	FindSLs(ctx context.Context, filter SLFilter) ([]*SuspiciousLogin, int, error)
	CreateSL(ctx context.Context, sl *SuspiciousLogin) error
	UpdateSL(ctx context.Context, id uint, upd SLUpdate) (*SuspiciousLogin, error)
	DeleteSL(ctx context.Context, id uint) error
}

type SLFilter struct {
	ID                 *uint   `json:"id"`
	UserID             *uint   `json:"user_id"`
	Email              *string `json:"email"`
	IP                 *string `json:"ip"`
	Country            *string `json:"country"`
	City               *string `json:"city"`
	Browser            *string `json:"browser"`
	Platform           *string `json:"platform"`
	OS                 *string `json:"os"`
	Device             *string `json:"device"`
	DeviceType         *string `json:"device_type"`
	UnverifiedAttempts *int    `json:"unverified_attempts"`
	IsTrusted          *bool   `json:"is_trusted"`
	IsBlocked          *bool   `json:"is_blocked"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SLUpdate struct {
	ID                 *uint   `json:"id"`
	UserID             *uint   `json:"user_id"`
	Email              *string `json:"email"`
	IP                 *string `json:"ip"`
	Country            *string `json:"country"`
	City               *string `json:"city"`
	Browser            *string `json:"browser"`
	Platform           *string `json:"platform"`
	OS                 *string `json:"os"`
	Device             *string `json:"device"`
	DeviceType         *string `json:"device_type"`
	UnverifiedAttempts *int    `json:"unverified_attempts"`
	IsTrusted          *bool   `json:"is_trusted"`
	IsBlocked          *bool   `json:"is_blocked"`
}
