package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

type TokenService struct {
	db *DB
}

func NewTokenService(db *DB) *TokenService {
	return &TokenService{db: db}
}

func (s *TokenService) FindTokenByID(ctx context.Context, id uint) (*sm.Token, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	token, err := findTokenByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenService) FindTokens(ctx context.Context, filter sm.TokenFilter) ([]*sm.Token, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findTokens(ctx, tx, filter)
}

func (s *TokenService) CreateToken(ctx context.Context, token *sm.Token) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createToken(ctx, tx, token); err != nil {
		return err
	}

	return tx.Commit()
}

func findTokenByID(ctx context.Context, tx *Tx, id uint) (*sm.Token, error) {
	a, _, err := findTokens(ctx, tx, sm.TokenFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Token is not found."}
	}
	return a[0], nil
}

func findTokens(ctx context.Context, tx *Tx, filter sm.TokenFilter) (_ []*sm.Token, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.UserID; v != nil {
		where, args = append(where, fmt.Sprintf(`"user_id" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "user_id", "refresh_token", "access_token", "created_at", "updated_at", COUNT(*) OVER()
	FROM "tokens"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	tokens := make([]*sm.Token, 0)
	for rows.Next() {
		var token sm.Token
		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.RefreshToken,
			&token.AccessToken,
			(*NullTime)(&token.CreatedAt),
			(*NullTime)(&token.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		tokens = append(tokens, &token)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tokens, n, nil
}

func createToken(ctx context.Context, tx *Tx, token *sm.Token) error {
	token.CreatedAt = tx.now
	token.UpdatedAt = token.CreatedAt

	query := `INSERT INTO "tokens" ("user_id", "refresh_token", "access_token", "created_at", "updated_at") 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	args := []interface{}{
		token.UserID,
		token.RefreshToken,
		token.AccessToken,
		(*NullTime)(&token.CreatedAt),
		(*NullTime)(&token.UpdatedAt),
	}

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&token.ID); err != nil {
		return err
	}

	return nil
}
