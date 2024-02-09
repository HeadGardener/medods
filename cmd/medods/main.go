package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/HeadGardener/medods/internal/config"
	"github.com/HeadGardener/medods/internal/handlers"
	"github.com/HeadGardener/medods/internal/lib/auth"
	"github.com/HeadGardener/medods/internal/server"
	"github.com/HeadGardener/medods/internal/services"
	"github.com/HeadGardener/medods/internal/storage"
)

const shutdownTimeout = 5 * time.Second

var confPath = flag.String("conf-path", "./config/.env", "path to config env")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	conf, err := config.Init(*confPath)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while initializing config: %s", err.Error())
	}

	db, err := storage.NewMongoCollection(ctx, &conf.DBConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %s", err.Error())
	}

	var (
		sessionStorage = storage.NewSessionStorage(db)
	)

	tokenManager, err := auth.NewTokenManager(&conf.TokensConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while initializing token manager: %s", err.Error())
	}

	var (
		authService = services.NewAuthService(tokenManager, sessionStorage)
	)

	handler := handlers.NewHandler(authService)

	srv := &server.Server{}
	go func() {
		if err = srv.Run(conf.ServerConfig, handler.InitRoutes()); err != nil {
			log.Printf("[ERROR] failed to run server: %s", err.Error())
		}
	}()
	log.Println("[INFO] server start working")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("[INFO] server forced to shutdown: %s", err.Error())
	}

	if err = db.Database().Client().Disconnect(ctx); err != nil {
		log.Printf("[INFO] db connection forced to shutdown: %s", err.Error())
	}

	log.Println("[INFO] server exiting")
}
