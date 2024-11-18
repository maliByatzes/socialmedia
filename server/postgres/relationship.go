package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

type RelationshipService struct {
	db *DB
}

func NewRelationshipService(db *DB) *RelationshipService {
	return &RelationshipService{db: db}
}

func (s *RelationshipService) FindRelationshipByID(ctx context.Context, id uint) (*sm.Relationship, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	r, err := findRelationshipByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *RelationshipService) FindRelationships(ctx context.Context, filter sm.RelationshipFilter) ([]*sm.Relationship, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findRelationships(ctx, tx, filter)
}

func (s *RelationshipService) CreateRelationship(ctx context.Context, r *sm.Relationship) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createRelationship(ctx, tx, r); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *RelationshipService) DeleteRelationship(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deleteRelationship(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func findRelationshipByID(ctx context.Context, tx *Tx, id uint) (*sm.Relationship, error) {
	a, _, err := findRelationships(ctx, tx, sm.RelationshipFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Relationship not found."}
	}
	return a[0], nil
}

func findRelationships(ctx context.Context, tx *Tx, filter sm.RelationshipFilter) (_ []*sm.Relationship, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.FollowerID; v != nil {
		where, args = append(where, fmt.Sprintf(`"follower_id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.FollowingID; v != nil {
		where, args = append(where, fmt.Sprintf(`"following_id" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "follower_id", "following_id", "created_at", "updated_at", COUNT(*) OVER()
	FROM "relationships"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	rs := make([]*sm.Relationship, 0)
	for rows.Next() {
		var r sm.Relationship
		if err := rows.Scan(
			&r.ID,
			&r.FollowerID,
			&r.FollowingID,
			(*NullTime)(&r.CreatedAt),
			(*NullTime)(&r.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		rs = append(rs, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return rs, n, nil
}

func createRelationship(ctx context.Context, tx *Tx, r *sm.Relationship) error {
	r.CreatedAt = tx.now
	r.UpdatedAt = r.CreatedAt

	query := `INSERT INTO "relationships" ("follower_id", "following_id", "created_at", "updated_at")
	VALUES ($1, $2, $3, $4) RETURNING id`
	args := []interface{}{
		r.FollowerID,
		r.FollowingID,
		(*NullTime)(&r.CreatedAt),
		(*NullTime)(&r.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

func deleteRelationship(ctx context.Context, tx *Tx, id uint) error {
	query := `DELETE FROM "relationships" WHERE id = $1`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
