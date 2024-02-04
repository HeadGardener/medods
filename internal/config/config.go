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
	SecretKey      string
	AccessTokenTTL time.Duration

	InitialLen      int
	RefreshTokenTTL time.Duration
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

	accessTokenSecretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if accessTokenSecretKey == "" {
		return nil, errors.New("access token secret key is empty")
	}

	accessTokenTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid access token ttl: %w", err)
	}

	refreshInitialLen, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_INITIAL_LEN"))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token initial len: %w", err)
	}

	refreshTokenTTL, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token ttl: %w", err)
	}

	return &Config{
		DBConfig: DBConfig{
			URL: dbURL,
		},
		ServerConfig: ServerConfig{
			Port: serverPort,
		},
		TokensConfig: TokensConfig{
			SecretKey:       accessTokenSecretKey,
			AccessTokenTTL:  time.Duration(accessTokenTTL) * time.Minute,
			InitialLen:      refreshInitialLen,
			RefreshTokenTTL: time.Duration(refreshTokenTTL) * time.Minute,
		},
	}, nil
}
