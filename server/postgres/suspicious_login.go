package postgres

import (
	"context"

	sm "github.com/maliByatzes/socialmedia"
)

type SuspiciousLoginService struct {
	db *DB
}

func NewSuspiciousLoginService(db *DB) *SuspiciousLoginService {
	return &SuspiciousLoginService{db: db}
}

func (s *SuspiciousLoginService) FindSLByID(ctx context.Context, id uint) (*sm.SuspiciousLogin, error) {
	return nil, nil
}

func (s *SuspiciousLoginService) FindSLByUserID(ctx context.Context, userID uint) (*sm.SuspiciousLogin, error) {
	return nil, nil
}

func (s *SuspiciousLoginService) FindSLs(ctx context.Context, filter sm.SLFilter) ([]*sm.SuspiciousLogin, int, error) {
	return nil, 0, nil
}

func (s *SuspiciousLoginService) CreateSL(ctx context.Context, sl *sm.SuspiciousLogin) error {
	return nil
}

func (s *SuspiciousLoginService) UpdateSL(ctx context.Context, id uint, upd sm.SLUpdate) (*sm.SuspiciousLogin, error) {
	return nil, nil
}

func (s *SuspiciousLoginService) DeleteSL(ctx context.Context, id uint) error {
	return nil
}
