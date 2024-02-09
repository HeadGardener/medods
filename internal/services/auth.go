package services

import (
	"context"
	"errors"
	"time"

	"github.com/HeadGardener/medods/internal/lib/auth"
	"github.com/HeadGardener/medods/internal/lib/hash"
	"github.com/HeadGardener/medods/internal/models"

	"github.com/google/uuid"
)

var (
	ErrNotSameRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidSession      = errors.New("tokens are not connected: invalid access token session")
)

type TokenManager interface {
	GenerateAccessToken(userID, sessionID string) (string, error)
	ParseAccessTokenWithoutExpirationTime(accessToken string) (auth.TokenAttributes, error)
	GenerateRefreshToken() (string, error)
	GetRefreshTokenTTL() time.Duration
}

type SessionStorage interface {
	GetSessionByUserID(ctx context.Context, userID string) (models.Session, error)
	CreateUserSession(ctx context.Context, session models.Session) error
	UpdateSession(ctx context.Context, session models.Session) error
	GetSessionByID(ctx context.Context, id string) (models.Session, error)
}

type AuthService struct {
	tokenManager   TokenManager
	sessionStorage SessionStorage
}

func NewAuthService(tokenManager TokenManager, sessionStorage SessionStorage) *AuthService {
	return &AuthService{
		tokenManager:   tokenManager,
		sessionStorage: sessionStorage,
	}
}

func (s *AuthService) SignIn(ctx context.Context, userID string) (models.Tokens, error) {
	if _, err := s.sessionStorage.GetSessionByUserID(ctx, userID); err != nil {
		session := models.Session{
			ID:           "",
			UserID:       userID,
			RefreshToken: "",
			ExpiresAt:    time.Now(),
		}

		if err = s.sessionStorage.CreateUserSession(ctx, session); err != nil {
			return models.Tokens{}, err
		}
	}

	tokens, err := s.createSession(ctx, userID)
	if err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}

func (s *AuthService) Refresh(ctx context.Context, accessToken, refreshToken string) (models.Tokens, error) {
	tokenAttr, err := s.tokenManager.ParseAccessTokenWithoutExpirationTime(accessToken)
	if err != nil {
		return models.Tokens{}, err
	}

	session, err := s.sessionStorage.GetSessionByID(ctx, tokenAttr.SessionID)
	if err != nil {
		return models.Tokens{}, err
	}

	if !hash.CompareHashAndString(session.RefreshToken, refreshToken) {
		return models.Tokens{}, ErrNotSameRefreshToken
	}

	if session.ExpiresAt.Before(time.Now()) {
		return models.Tokens{}, ErrRefreshTokenExpired
	}

	tokens, err := s.createSession(ctx, session.UserID)
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

	sessionID := uuid.NewString()
	tokens.AccessToken, err = s.tokenManager.GenerateAccessToken(userID, sessionID)
	if err != nil {
		return models.Tokens{}, err
	}

	tokens.RefreshToken, err = s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return models.Tokens{}, err
	}

	session := models.Session{
		ID:           sessionID,
		UserID:       userID,
		RefreshToken: hash.GetStringHash(tokens.RefreshToken),
		ExpiresAt:    time.Now().Add(s.tokenManager.GetRefreshTokenTTL()),
	}

	if err = s.sessionStorage.UpdateSession(ctx, session); err != nil {
		return models.Tokens{}, err
	}

	return tokens, nil
}
