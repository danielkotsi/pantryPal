package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"pantrypal/backend/internal/platform/auth"
	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

var (
	ErrInvalidEmail       = errors.New("email is invalid")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrEmailConflict      = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	users  *repositories.UserRepository
	tokens *auth.TokenManager
}

func NewAuthService(users *repositories.UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if !strings.Contains(req.Email, "@") {
		return dto.AuthResponse{}, ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return dto.AuthResponse{}, ErrWeakPassword
	}

	exists, err := s.users.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if exists {
		return dto.AuthResponse{}, ErrEmailConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	user, err := s.users.CreateUser(ctx, req.Email, string(hash), req.DisplayName)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	token, expiresAt, err := s.tokens.MakeToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.users.GetByEmail(ctx, req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.AuthResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return dto.AuthResponse{}, ErrInvalidCredentials
	}

	if err := s.users.EnsureDefaultProfileRows(ctx, user.ID); err != nil {
		return dto.AuthResponse{}, err
	}

	token, expiresAt, err := s.tokens.MakeToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (dto.UserResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{ID: user.ID, Email: user.Email, DisplayName: user.DisplayName}, nil
}
