package auth

import (
	"time"

	"github.com/HeadGardener/medods/internal/config"
)

type TokenProcessor struct {
	accessTokenParams  tokenParams
	refreshTokenParams tokenParams
}

type tokenParams struct {
	PublicKey  string
	PrivateKey string
	ExpiresIn  time.Duration
}

func NewTokenProcessor(conf *config.TokensConfig) *TokenProcessor {
	return &TokenProcessor{
		accessTokenParams: tokenParams{
			PublicKey:  conf.AccessToken.PublicKey,
			PrivateKey: conf.AccessToken.PrivateKey,
			ExpiresIn:  conf.AccessToken.ExpiresIn,
		},
		refreshTokenParams: tokenParams{
			PublicKey:  conf.RefreshToken.PublicKey,
			PrivateKey: conf.RefreshToken.PrivateKey,
			ExpiresIn:  conf.RefreshToken.ExpiresIn,
		},
	}
}
