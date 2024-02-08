package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const (
	userIDQueryParam = "user_id"
)

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get(userIDQueryParam)
	if userID == "" {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id from url",
			errors.New("user_id is empty"))
		return
	}

	if _, err := uuid.Parse(userID); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating user id", err)
		return
	}

	tokens, err := h.authService.SignIn(r.Context(), userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while signing in", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
