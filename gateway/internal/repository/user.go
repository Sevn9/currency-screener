package repository

import (
	"context"
	"sync"

	"github.com/Sevn9/currency-screener/gateway/internal/models"
	"github.com/Sevn9/currency-screener/pkg/apperrors"
)

type User = models.User

type UserRepository struct {
	users map[string]User
	mu    *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]User),
		mu:    &sync.RWMutex{},
	}
}

func (repo *UserRepository) AddUser(user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.Login]; exists {
		return apperrors.ErrAlreadyExists
	}
	repo.users[user.Login] = user

	return nil
}

func (repo *UserRepository) GetUser(ctx context.Context, login string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[login]
	if !exists {
		return User{}, apperrors.NotFoundError{Entity: "User", ID: login}
	}

	return user, nil
}
