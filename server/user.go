package socialmedia

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	Avatar          string    `json:"avatar"`
	Location        string    `json:"location"`
	Bio             string    `json:"bio"`
	Interests       string    `json:"interests"`
	Role            string    `json:"role"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("user: name is required.")
	}

	if u.Email == "" {
		return fmt.Errorf("user: email is required.")
	}

	if u.Password == "" {
		return fmt.Errorf("user: password is required.")
	}

	return nil
}

type UserService interface {
  CreateUser(ctx context.Context, user *User) error
}
