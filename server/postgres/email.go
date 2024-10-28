package postgres

import (
	"context"
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

func (s *EmailService) CreateEmailVerification(ctx context.Context, email *sm.Email) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createEmailVerification(ctx, tx, email); err != nil {
		return err
	}

	return tx.Commit()
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
