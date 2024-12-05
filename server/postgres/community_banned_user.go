package postgres

import (
	"context"

	sm "github.com/maliByatzes/socialmedia"
)

type CBUService struct {
	db *DB
}

func NewCBUService(db *DB) *CBUService {
	return &CBUService{db: db}
}

func (s *CBUService) FindCBUs(ctx context.Context, filter sm.CBUFilter) ([]*sm.CBU, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findCBUs(ctx, tx, filter)
}

func (s *CBUService) CreateCBU(ctx context.Context, cbu *sm.CBU) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createCBU(ctx, tx, cbu); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CBUService) DeleteCBU(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteCBU(ctx, tx, id); err != nil {
		return err
	} else if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func findCBUs(ctx context.Context, tx *Tx, filter sm.CBUFilter) (_ []*sm.CBU, n int, err error) {

}

func createCBU(ctx context.Context, tx *Tx, cbu *sm.CBU) error {

}

func deleteCBU(ctx context.Context, tx *Tx, id uint) error {

}
