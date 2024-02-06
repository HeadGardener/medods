package services

import (
	"context"
	"errors"
	"time"

	"github.com/HeadGardener/medods/internal/lib/hash"
	"github.com/HeadGardener/medods/internal/models"
)

var (
	ErrNotSameRefreshToken = errors.New("invalid refresh token")
)

type TokenManager interface {
	GenerateAccessToken(userID string) (string, error)
	ParseAccessToken(accessToken string) (string, error)
	GenerateRefreshToken() (string, error)
	GetRefreshTokenTTL() time.Duration
}

type SessionStorage interface {
	IsUserSessionExists(ctx context.Context, userID string) bool
	CreateUserSession(ctx context.Context, session models.Session) error
	UpdateSession(ctx context.Context, session models.Session) error
	GetSession(ctx context.Context, userID string) (models.Session, error)
}

type AuthService struct {
	tokenManager   TokenManager
	sessionStorage SessionStorage
}

func (s *AuthService) SignIn(ctx context.Context, userID string) (models.Tokens, error) {
	if !s.sessionStorage.IsUserSessionExists(ctx, userID) {
		session := models.Session{
			UserID:       userID,
			RefreshToken: "",
			ExpiresAt:    time.Now(),
		}

		if err := s.sessionStorage.CreateUserSession(ctx, session); err != nil {
			return models.Tokens{}, err
		}
	}

	tokens, err := s.createSession(ctx, userID)
	if err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, userID, refreshToken string) (models.Tokens, error) {
	session, err := s.sessionStorage.GetSession(ctx, userID)
	if err != nil {
		return models.Tokens{}, err
	}

	if !hash.CompareHashAndString(session.RefreshToken, refreshToken) {
		return models.Tokens{}, ErrNotSameRefreshToken
	}

	tokens, err := s.createSession(ctx, userID)
	if err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}

func (s *AuthService) createSession(ctx context.Context, userID string) (models.Tokens, error) {
	var (
		tokens models.Tokens
		err    error
	)

	tokens.AccessToken, err = s.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens.RefreshToken, err = s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return models.Tokens{}, err
	}

	session := models.Session{
		UserID:       userID,
		RefreshToken: hash.GetStringHash(tokens.RefreshToken),
		ExpiresAt:    time.Now().Add(s.tokenManager.GetRefreshTokenTTL()),
	}

	if err = s.sessionStorage.UpdateSession(ctx, session); err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}
