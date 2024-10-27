package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

var _ sm.UserService = (*UserService)(nil)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, user *sm.User) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func createUser(ctx context.Context, tx *Tx, user *sm.User) error {
	user.CreatedAt = tx.now
	user.UpdatedAt = user.CreatedAt

	if err := user.Validate(); err != nil {
		return err
	}

	query := `
  INSERT INTO "users" ("name", "email", "password", "avatar", "role", "created_at", "updated_at")
  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
  `
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password,
		user.Avatar,
		user.Role,
		(*NullTime)(&user.CreatedAt),
		(*NullTime)(&user.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return sm.Errorf(sm.ECONFLICT, "this email is already exists.")
		default:
			return err
		}
	}

	return nil
}
