package postgres

import (
	"context"
	"fmt"

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

func findCommunityByID(ctx context.Context, tx *Tx, id uint) (*sm.Community, error) {
	a, _, err := findCommunities(ctx, tx, sm.CommunityFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Community not found."}
	}
	return a[0], nil
}

func findCommunities(ctx context.Context, tx *Tx, filter sm.CommunityFilter) (_ []*sm.Community, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.Name; v != nil {
		where, args = append(where, fmt.Sprintf(`"name" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "name", "description", "banner", "created_at", "updated_at", COUNT(*) OVER()
	FROM "communities"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	coms := make([]*sm.Community, 0)
	for rows.Next() {
		var com sm.Community
		if err := rows.Scan(
			&com.ID,
			&com.Name,
			(*NullString)(&com.Description),
			(*NullString)(&com.Banner),
			(*NullTime)(&com.CreatedAt),
			(*NullTime)(&com.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		coms = append(coms, &com)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return coms, n, nil
}

func createCommunity(ctx context.Context, tx *Tx, com *sm.Community) error {
	com.CreatedAt = tx.now
	com.UpdatedAt = com.CreatedAt

	query := `INSERT INTO "communities"("name", "description", "banner", "created_at", "updated_at")
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	args := []interface{}{
		com.Name,
		com.Description,
		com.Banner,
		com.CreatedAt,
		com.UpdatedAt,
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&com.ID)
	if err != nil {
		return err
	}

	return nil
}

func updateCommunity(ctx context.Context, tx *Tx, id uint, upd sm.CommunityUpdate) (*sm.Community, error) {
	com, err := findCommunityByID(ctx, tx, id)
	if err != nil {
		return com, err
	}

	if v := upd.Name; v != nil {
		com.Name = *v
	}

	if v := upd.Description; v != nil {
		com.Description = *v
	}

	if v := upd.Banner; v != nil {
		com.Banner = *v
	}

	com.UpdatedAt = tx.now

	args := []interface{}{
		com.Name,
		com.Description,
		com.Banner,
		com.UpdatedAt,
		com.ID,
	}
	query := `UPDATE "communities" SET "name" = $1, "description" = $2, "banner" = $3, "updated_at" = $4 WHERE "id" = $5`

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return com, err
	}

	return com, err
}

func deleteCommunity(ctx context.Context, tx *Tx, id uint) error {
	_, err := findCommunityByID(ctx, tx, id)
	if err != nil {
		return err
	}

	query := `DELETE FROM "communities" WHERE "id" = $1`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
