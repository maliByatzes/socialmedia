package postgres

import (
	"context"

	sm "github.com/maliByatzes/socialmedia"
)

var _ sm.PreferenceService = (*PreferenceService)(nil)

type PreferenceService struct {
	db *DB
}

func NewPreferenceService(db *DB) *PreferenceService {
	return &PreferenceService{db: db}
}
func (s *PreferenceService) FindPreferenceByID(ctx context.Context, id uint) (*sm.Preference, error) {
	panic("not implemented")
}

func (s *PreferenceService) FindPreferenceByUserID(ctx context.Context, userID uint) (*sm.Preference, error) {
	panic("not implemented")
}

func (s *PreferenceService) FindPreferences(ctx context.Context, filter sm.PreferenceFilter) ([]*sm.Preference, int, error) {
	panic("not implemented")
}

func (s *PreferenceService) CreatePreference(ctx context.Context, preference *sm.Preference) error {
	panic("not implemented")
}

func (s *PreferenceService) UpdatePreference(ctx context.Context, id uint, upd sm.PreferenceUpdate) (*sm.Preference, error) {
	panic("not implemented")
}

func (s *PreferenceService) DeletePreference(ctx context.Context, id uint) error {
	panic("not implemented")
}
