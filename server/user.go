package socialmedia

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
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
		return Errorf(EINVALID, "name is required.")
	}

	if u.Email == "" {
		return Errorf(EINVALID, "email is required.")
	}

	if u.Password == "" {
		return Errorf(EINVALID, "password is required.")
	}

	return nil
}

func (u *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashBytes)

	return nil
}

func (u *User) VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
}
