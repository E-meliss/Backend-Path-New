package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func (u *User) ValidateForCreate(plainPassword string) error {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 chars")
	}
	if !emailRe.MatchString(u.Email) {
		return errors.New("invalid email")
	}
	if len(plainPassword) < 8 {
		return errors.New("password must be at least 8 chars")
	}
	return nil
}
