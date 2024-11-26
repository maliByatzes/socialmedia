package postgres

import (
	"context"
	"fmt"

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
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.CommunityID; v != nil {
		where, args = append(where, fmt.Sprintf(`"community_id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.UserID; v != nil {
		where, args = append(where, fmt.Sprintf(`"user_id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.IsModerator; v != nil {
		where, args = append(where, fmt.Sprintf(`"is_moderator" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "community_id", "user_id", "is_moderator", "created_at", "updated_at", COUNT(*) OVER()
	FROM "community_members"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	cms := make([]*sm.CommunityMember, 0)
	for rows.Next() {
		var cm sm.CommunityMember
		if err := rows.Scan(
			&cm.CommunityID,
			&cm.UserID,
			&cm.IsModerator,
			(*NullTime)(&cm.CreatedAt),
			(*NullTime)(&cm.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		cms = append(cms, &cm)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return cms, n, nil
}

func createCommunityMember(ctx context.Context, tx *Tx, cm *sm.CommunityMember) error {
	cm.CreatedAt = tx.now
	cm.UpdatedAt = cm.CreatedAt

	query := `INSERT INTO "community_members"("community_id", "user_id", "is_moderator", "created_at", "updated_at")
	VALUES ($1, $2, $3, $4, $5)`
	args := []interface{}{
		cm.CommunityID,
		cm.UserID,
		cm.IsModerator,
		(*NullTime)(&cm.CreatedAt),
		(*NullTime)(&cm.UpdatedAt),
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func updateCommunityMember(ctx context.Context, tx *Tx, id uint, upd sm.CommunityMemberUpdate) (*sm.CommunityMember, error) {
}

func deleteCommunityMember(cxt context.Context, tx *Tx, id uint) error {}
