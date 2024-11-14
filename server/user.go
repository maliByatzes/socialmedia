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
	Password        string    `json:"-" db:"password"`
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

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type UserService interface {
	FindUserByID(ctx context.Context, id uint) (*User, error)
	Authenticate(ctx context.Context, email, password string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)
	UpdateUser(ctx context.Context, id uint, up UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type UserFilter struct {
	ID    *uint   `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type UserUpdate struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Avatar          *string `json:"avatar"`
	Location        *string `json:"location"`
	Bio             *string `json:"bio"`
	Interests       *string `json:"interests"`
	Role            *string `json:"role"`
	IsEmailVerified *bool   `json:"is_email_verified"`
}
