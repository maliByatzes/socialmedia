package postgres

import (
	"context"
	"fmt"
	"strings"

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

func (s *UserService) FindUserByID(ctx context.Context, id uint) (*sm.User, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*sm.User, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	user, err := findUserByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}

	if err := user.VerifyPassword(password); err != nil {
		return nil, &sm.Error{Code: sm.ENOTAUTHORIZED, Message: "Invalid credentials"}
	}

	return user, nil
}

func (s *UserService) FindUsers(ctx context.Context, filter sm.UserFilter) ([]*sm.User, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findUsers(ctx, tx, filter)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, up sm.UserUpdate) (*sm.User, error) {
	return nil, sm.Errorf(sm.ENOTIMPLEMENTED, "")
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return sm.Errorf(sm.ENOTIMPLEMENTED, "")
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

func findUserByID(ctx context.Context, tx *Tx, id uint) (*sm.User, error) {
	a, _, err := findUsers(ctx, tx, sm.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "User not found"}
	}
	return a[0], nil
}

func findUserByEmail(ctx context.Context, tx *Tx, email string) (*sm.User, error) {
	a, _, err := findUsers(ctx, tx, sm.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "User not found"}
	}
	return a[0], nil
}

func findUsers(ctx context.Context, tx *Tx, filter sm.UserFilter) (_ []*sm.User, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPosition)), append(args, *v)
	}

	if v := filter.Name; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf(`"name" = $%d`, argPosition)), append(args, *v)
	}

	if v := filter.Email; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf(`"email" = $%d`, argPosition)), append(args, *v)
	}

	query := `SELECT "id", "name", "email", "password", "avatar", "location",
  "bio", "interests", "role", "is_email_verified", "created_at", "updated_at",
  COUNT(*) OVER() FROM "users"` + formatWhereClause(where) + ` ORDER BY id
  ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	users := make([]*sm.User, 0)
	for rows.Next() {
		var user sm.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			(*NullString)(&user.Avatar),
			(*NullString)(&user.Location),
			(*NullString)(&user.Bio),
			(*NullString)(&user.Interests),
			(*NullString)(&user.Role),
			&user.IsEmailVerified,
			(*NullTime)(&user.CreatedAt),
			(*NullTime)(&user.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, n, nil
}

func updateUser(ctx context.Context, tx *Tx, id uint, upd sm.UserUpdate) (*sm.User, error) {
	panic("Not implemented")
}

func deleteUser(ctx context.Context, tx *Tx, id uint) error {
	panic("Not implemented")
}

func formatWhereClause(where []string) string {
	if len(where) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(where, " AND ")
}

func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}
