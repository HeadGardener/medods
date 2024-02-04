package auth

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"time"

	"github.com/HeadGardener/medods/internal/models"
	"github.com/golang-jwt/jwt/v5"

	"github.com/HeadGardener/medods/internal/config"
)

type TokenProcessor struct {
	SecretKey       []byte
	AccessTokenTTL  time.Duration
	InitialLen      int
	RefreshTokenTTL time.Duration
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func NewTokenProcessor(conf *config.TokensConfig) (*TokenProcessor, error) {
	return &TokenProcessor{
		SecretKey:       []byte(conf.SecretKey),
		AccessTokenTTL:  conf.AccessTokenTTL,
		InitialLen:      conf.InitialLen,
		RefreshTokenTTL: conf.RefreshTokenTTL,
	}, nil
}

func (tp *TokenProcessor) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tp.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
	})

	return token.SignedString(tp.SecretKey)
}

func (tp *TokenProcessor) ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return tp.SecretKey, nil
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

func (tp *TokenProcessor) GenerateRefreshToken() (models.Session, error) {
	b := make([]byte, tp.InitialLen)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	_, err := r.Read(b)
	if err != nil {
		return models.Session{}, err
	}

	token := base64.StdEncoding.EncodeToString(b)

	return models.Session{
		RefreshToken: token,
		ExpiresAt:    time.Now().Add(tp.RefreshTokenTTL),
	}, nil
}
