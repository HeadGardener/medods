package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	TokensConfig TokensConfig
}

type DBConfig struct {
	DBName     string
	Collection string
	URL        string
	Username   string
	Password   string
}

type ServerConfig struct {
	Port string
}

type TokensConfig struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

type AccessToken struct {
	PublicKey  string
	PrivateKey string
	ExpiresIn  time.Duration
}

type RefreshToken struct {
	PublicKey  string
	PrivateKey string
	ExpiresIn  time.Duration
}

func Init(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		return nil, errors.New("db name is empty")
	}

	coll := os.Getenv("COLLECTION")
	if coll == "" {
		return nil, errors.New("collection is empty")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("db url is empty")
	}

	dbUsername := os.Getenv("DATABASE_USER")
	if dbUsername == "" {
		return nil, errors.New("db username is empty")
	}

	dbPassword := os.Getenv("DATABASE_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("db password is empty")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return nil, errors.New("server port is empty")
	}

	accessTokenPublicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	if accessTokenPublicKey == "" {
		return nil, errors.New("access token public key is empty")
	}

	accessTokenPrivateKey := os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	if accessTokenPrivateKey == "" {
		return nil, errors.New("access token private key is empty")
	}

	accessTokenExpiresIn, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRES_IN"))
	if err != nil {
		return nil, fmt.Errorf("invalid access token expires in env: %w", err)
	}

	refreshTokenPublicKey := os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	if refreshTokenPublicKey == "" {
		return nil, errors.New("access token public key is empty")
	}

	refreshTokenPrivateKey := os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	if refreshTokenPrivateKey == "" {
		return nil, errors.New("access token private key is empty")
	}

	refreshTokenExpiresIn, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRES_IN"))
	if err != nil {
		return nil, fmt.Errorf("invalid access token expires in env: %w", err)
	}

	return &Config{
		DBConfig: DBConfig{
			URL: dbURL,
		},
		ServerConfig: ServerConfig{
			Port: serverPort,
		},
		TokensConfig: TokensConfig{
			AccessToken: AccessToken{
				PublicKey:  accessTokenPublicKey,
				PrivateKey: accessTokenPrivateKey,
				ExpiresIn:  time.Duration(accessTokenExpiresIn) * time.Minute,
			},
			RefreshToken: RefreshToken{
				PublicKey:  refreshTokenPublicKey,
				PrivateKey: refreshTokenPrivateKey,
				ExpiresIn:  time.Duration(refreshTokenExpiresIn) * time.Minute,
			},
		},
	}, nil
}
