package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/E-meliss/wallet-service/internal/domain"
)

type UserService struct {
	users domain.UserRepository
}

func NewUserService(users domain.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (domain.User, error) {
	u := domain.User{Username: username, Email: email, Role: domain.RoleUser}
	if err := u.ValidateForCreate(password); err != nil {
		return domain.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	u.PasswordHash = string(hash)

	return s.users.Create(ctx, u)
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (domain.User, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return domain.User{}, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *UserService) Authorize(u domain.User, required domain.UserRole) error {
	if u.Role == domain.RoleAdmin {
		return nil
	}
	if u.Role != required {
		return errors.New("forbidden")
	}
	return nil
}
