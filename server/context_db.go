package socialmedia

import (
	"context"
	"time"
)

type Context struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Email      string    `json:"email"`
	IP         string    `json:"ip"`
	Country    string    `json:"country"`
	City       string    `json:"city"`
	Browser    string    `json:"browser"`
	Platform   string    `json:"platform"`
	OS         string    `json:"os"`
	Device     string    `json:"device"`
	DeviceType string    `json:"device_type"`
	IsTrusted  bool      `json:"is_trusted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Context) Validate() error {
	if c.UserID == 0 {
		return Errorf(EINVALID, "UserID is required.")
	}

	return nil
}

type ContextService interface {
	FindContexts(ctx context.Context, filter ContextFilter) ([]*Context, int, error)
	FindContextByID(ctx context.Context, id uint) (*Context, error)
	FindContextByUserID(ctx context.Context, userID uint) (*Context, error)
	CreateContext(ctx context.Context, context *Context) error
	UpdateContext(ctx context.Context, id uint, upd ContextUpdate) (*Context, error)
	DeleteContext(ctx context.Context, id uint) error
}

type ContextFilter struct {
	ID         *uint   `json:"id"`
	UserID     *uint   `json:"user_id"`
	Email      *string `json:"email"`
	IP         *string `json:"ip"`
	Country    *string `json:"country"`
	City       *string `json:"city"`
	Browser    *string `json:"browser"`
	Platform   *string `json:"platform"`
	OS         *string `json:"os"`
	Device     *string `json:"device"`
	DeviceType *string `json:"device_type"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ContextUpdate struct {
	Email      *string `json:"email"`
	IP         *string `json:"ip"`
	Country    *string `json:"country"`
	City       *string `json:"city"`
	Browser    *string `json:"browser"`
	Platform   *string `json:"platform"`
	OS         *string `json:"os"`
	Device     *string `json:"device"`
	DeviceType *string `json:"device_type"`
}
