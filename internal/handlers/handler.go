package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/HeadGardener/medods/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	minute = time.Minute
)

type AuthService interface {
	SignIn(ctx context.Context, userID string) (models.Tokens, error)
	Refresh(ctx context.Context, accessToken, refreshToken string) (models.Tokens, error)
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(minute))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-in", h.signIn)
			r.Route("/token", func(r chi.Router) {
				r.Use(h.extractToken)
				r.Put("/refresh", h.refresh)
			})
		})
	})

	return r
}
