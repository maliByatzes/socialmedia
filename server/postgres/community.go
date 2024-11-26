package postgres

import (
	"context"

	sm "github.com/maliByatzes/socialmedia"
)

type CommunityService struct {
	db *DB
}

func NewCommunityService(db *DB) *CommunityService {
	return &CommunityService{db: db}
}

func (s *CommunityService) FindCommunityByID(ctx context.Context, id uint) (*sm.Community, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	com, err := findCommunityByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return com, nil
}

func (s *CommunityService) FindCommunities(ctx context.Context, filter sm.CommunityFilter) ([]*sm.Community, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findCommunities(ctx, tx, filter)
}

func (s *CommunityService) CreateCommunity(ctx context.Context, com *sm.Community) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createCommunity(ctx, tx, com); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CommunityService) UpdateCommunity(ctx context.Context, id uint, upd sm.CommunityUpdate) (*sm.Community, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	com, err := updateCommunity(ctx, tx, id, upd)
	if err != nil {
		return com, err
	} else if err := tx.Commit(); err != nil {
		return com, err
	}

	return com, nil
}

func (s *CommunityService) DeleteCommunity(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteCommunity(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func findCommunityByID(ctx context.Context, tx *Tx, id uint) (*sm.Community, error) {}

func findCommunities(ctx context.Context, tx *Tx, filter sm.CommunityFilter) (_ []*sm.Community, n int, err error) {
}

func createCommunity(ctx context.Context, tx *Tx, com *sm.Community) error {}

func updateCommunity(ctx context.Context, tx *Tx, id uint, upd sm.CommunityUpdate) (*sm.Community, error) {
}

func deleteCommunity(ctx context.Context, tx *Tx, id uint) error {}
