package postgres

import (
	"context"

	sm "github.com/maliByatzes/socialmedia"
)

type CommunityMemberService struct {
	db *DB
}

func NewCommunityMemberService(db *DB) *CommunityMemberService {
	return &CommunityMemberService{db: db}
}

func (s *CommunityMemberService) FindCommunityMembers(ctx context.Context, filter sm.CommunityMemberFilter) ([]*sm.CommunityMember, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findCommunityMembers(ctx, tx, filter)
}

func (s *CommunityMemberService) CreateCommunityMember(ctx context.Context, cm *sm.CommunityMember) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createCommunityMember(ctx, tx, cm); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CommunityMemberService) UpdateCommunityMember(ctx context.Context, id uint, upd sm.CommunityMemberUpdate) (*sm.CommunityMember, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	cm, err := updateCommunityMember(ctx, tx, id, upd)
	if err != nil {
		return cm, err
	} else if err := tx.Commit(); err != nil {
		return cm, err
	}

	return cm, nil
}

func (s *CommunityMemberService) DeleteCommunityMember(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteCommunityMember(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func findCommunityMembers(ctx context.Context, tx *Tx, filter sm.CommunityMemberFilter) (_ []*sm.CommunityMember, n int, err error) {

}

func createCommunityMember(ctx context.Context, tx *Tx, cm *sm.CommunityMember) error {}

func updateCommunityMember(ctx context.Context, tx *Tx, id uint, upd sm.CommunityMemberUpdate) (*sm.CommunityMember, error) {
}

func deleteCommunityMember(cxt context.Context, tx *Tx, id uint) error {}
