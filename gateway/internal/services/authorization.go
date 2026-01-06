package services

import (
	"context"
	"fmt"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/Sevn9/currency-screener/gateway/internal/models"
	"github.com/Sevn9/currency-screener/gateway/internal/repository"
	"github.com/Sevn9/currency-screener/pkg/apperrors"
)

type authClientInterface interface {
	GenerateToken(ctx context.Context, login string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}

type repositoryInterface interface {
	AddUser(user models.User) error
	GetUser(ctx context.Context, login string) (models.User, error)
}

type AuthService struct {
	authClient authClientInterface
	repository repositoryInterface
}

func NewAuth(authClient authClientInterface, repository repositoryInterface) AuthService {
	return AuthService{
		authClient: authClient,
		repository: repository,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) error {
	user := repository.User{Login: req.Username, Password: req.Password}
	if err := s.repository.AddUser(user); err != nil {
		return fmt.Errorf("services: repository.AddUser: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repository.GetUser(ctx, login)
	if err != nil {
		return "", fmt.Errorf("services: repository.GetUser: %w", err)
	}

	if user.Password != password {
		return "", apperrors.ErrInvalidCredentials
	}

	res, err := s.authClient.GenerateToken(ctx, login)
	if err != nil {
		return "", fmt.Errorf("services: authClient.GenerateToken: %w", err)
	}

	return res, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) error {
	err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return fmt.Errorf("services: authClient.ValidateToken: %w", err)
	}

	return nil
}

func (s *AuthService) Logout(token string) error {
	return apperrors.ErrNotImplemented
}
