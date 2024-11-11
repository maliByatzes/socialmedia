package token

import "time"

type Payload struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(id uint, name string, duration time.Duration) *Payload {
	return &Payload{
		ID:        id,
		Name:      name,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
