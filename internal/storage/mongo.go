package storage

import (
	"context"
	"time"

	"github.com/HeadGardener/medods/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connTimeout = 5 * time.Second
)

func NewMongoCollection(ctx context.Context, conf *config.DBConfig) (*mongo.Collection, error) {
	ctxConn, cancel := context.WithTimeout(ctx, connTimeout)
	defer cancel()

	credential := options.Credential{
		Username: conf.Username,
		Password: conf.Password,
	}
	clientOpts := options.Client().ApplyURI(conf.URL).SetAuth(credential)

	client, err := mongo.Connect(ctxConn, clientOpts)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, connTimeout)
	defer cancel()

	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(conf.DBName).Collection(conf.Collection), err
}
