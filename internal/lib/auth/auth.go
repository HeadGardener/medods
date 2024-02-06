package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/HeadGardener/medods/internal/config"
)

type TokenManager struct {
	SecretKey       []byte
	AccessTokenTTL  time.Duration
	InitialLen      int
	RefreshTokenTTL time.Duration
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func NewTokenProcessor(conf *config.TokensConfig) (*TokenManager, error) {
	return &TokenManager{
		SecretKey:       []byte(conf.SecretKey),
		AccessTokenTTL:  conf.AccessTokenTTL,
		InitialLen:      conf.InitialLen,
		RefreshTokenTTL: conf.RefreshTokenTTL,
	}, nil
}

func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
	})

	return token.SignedString(tm.SecretKey)
}

func (tm *TokenManager) ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return tm.SecretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, tm.InitialLen)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(b)

	return token, nil
}

func (tm *TokenManager) GetRefreshTokenTTL() time.Duration {
	return tm.RefreshTokenTTL
}
