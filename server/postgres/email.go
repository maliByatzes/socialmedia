package postgres

import (
	"context"
	"fmt"
	"time"

	sm "github.com/maliByatzes/socialmedia"
)

var _ sm.EmailService = (*EmailService)(nil)

type EmailService struct {
	db *DB
}

func NewEmailService(db *DB) *EmailService {
	return &EmailService{db: db}
}

func (s *EmailService) FindEmailVerificationByID(ctx context.Context, id uint) (*sm.Email, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	ev, err := findEmailVerificationByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *EmailService) FindEmailVerifications(ctx context.Context, filter sm.EmailFilter) ([]*sm.Email, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findEmailVerifications(ctx, tx, filter)
}

func (s *EmailService) CreateEmailVerification(ctx context.Context, email *sm.Email) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createEmailVerification(ctx, tx, email); err != nil {
		return err
	}

	return tx.Commit()
}

func findEmailVerificationByID(ctx context.Context, tx *Tx, id uint) (*sm.Email, error) {
	a, _, err := findEmailVerifications(ctx, tx, sm.EmailFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Email verification not found."}
	}
	return a[0], nil
}

func findEmailVerifications(ctx context.Context, tx *Tx, filter sm.EmailFilter) (_ []*sm.Email, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.Email; v != nil {
		where, args = append(where, fmt.Sprintf(`"email" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.VerificationCode; v != nil {
		where, args = append(where, fmt.Sprintf(`"verification_code" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.For; v != nil {
		where, args = append(where, fmt.Sprintf(`"for" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "email", "verification_code", "message_id", "for", "created_at", "expired_at", COUNT(*) OVER()
	FROM "emails"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}

	emails := make([]*sm.Email, 0)
	for rows.Next() {
		var email sm.Email
		if err := rows.Scan(
			&email.ID,
			&email.Email,
			&email.VerificationCode,
			&email.MessageID,
			&email.For,
			(*NullTime)(&email.CreatedAt),
			(*NullTime)(&email.ExpiresAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		emails = append(emails, &email)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return emails, n, nil
}

func createEmailVerification(ctx context.Context, tx *Tx, email *sm.Email) error {
	email.CreatedAt = tx.now
	email.ExpiresAt = email.CreatedAt.Add(30 * time.Minute)

	query := `
  INSERT INTO "emails" ("email", "verification_code", "message_id", "for", "created_at", "expires_at")
  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
  `
	args := []interface{}{
		email.Email,
		email.VerificationCode,
		email.MessageID,
		email.For,
		(*NullTime)(&email.CreatedAt),
		(*NullTime)(&email.ExpiresAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&email.ID)
	if err != nil {
		return err
	}

	return nil
}
