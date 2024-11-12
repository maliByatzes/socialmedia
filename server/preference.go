package socialmedia

import (
	"context"
	"time"
)

type Preference struct {
	ID                      uint      `json:"id"`
	UserID                  uint      `json:"user_id"`
	EnabledContextBasedAuth bool      `json:"enable_context_based_auth"`
	CreatedAt               time.Time `json:"created_at"`
}

func (p *Preference) Validate() error {
	if p.UserID == uint(0) {
		return Errorf(EINVALID, "UserID is required.")
	}

	return nil
}

type PreferenceService interface {
	FindPreferenceByID(ctx context.Context, id uint) (*Preference, error)
	FindPreferenceByUserID(ctx context.Context, userID uint) (*Preference, error)
	FindPreferences(ctx context.Context, filter PreferenceFilter) ([]*Preference, int, error)
	CreatePreference(ctx context.Context, preference *Preference) error
	UpdatePreference(ctx context.Context, id uint, upd PreferenceUpdate) (*Preference, error)
	DeletePreference(ctx context.Context, id uint) error
}

type PreferenceFilter struct {
	ID                      *uint `json:"id"`
	UserID                  *uint `json:"user_id"`
	EnabledContextBasedAuth *bool `json:"enable_context_based_auth"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type PreferenceUpdate struct {
	EnabledContextBasedAuth *bool `json:"enable_context_based_auth"`
}
