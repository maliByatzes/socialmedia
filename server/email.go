package socialmedia

import (
	"context"
	"time"
)

type Email struct {
	ID               uint      `json:"id"`
	Email            string    `json:"email"`
	VerificationCode string    `json:"verification_code"`
	MessageID        string    `json:"message_id"`
	For              string    `json:"for"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type EmailService interface {
	CreateEmailVerification(ctx context.Context, email *Email) error
}
