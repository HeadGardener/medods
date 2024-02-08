package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	tokenCtx = "access_token"
)

const (
	headerPartsLen = 2
)

var (
	ErrTokenNotString = errors.New("tokenCtx value is not of type string")
)

func (h *Handler) extractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("empty auth header"))
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != headerPartsLen {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("invalid auth header, must be like `Bearer token`"))
			return
		}

		if headerParts[0] != "Bearer" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				fmt.Errorf("invalid auth header %s, must be Bearer", headerParts[0]))
			return
		}

		token := headerParts[1]
		if token == "" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("jwt token is empty"))
			return
		}

		ctx := context.WithValue(r.Context(), tokenCtx, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getToken(r *http.Request) (string, error) {
	tokenCtxValue := r.Context().Value(tokenCtx)
	accessToken, ok := tokenCtxValue.(string)
	if !ok {
		return "", ErrTokenNotString
	}

	return accessToken, nil
}
