package handlers

import (
	"encoding/json"
	"net/http"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	accessToken, err := getToken(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while extracting access token", err)
		return
	}

	var req RefreshRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding RefreshRequest", err)
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), accessToken, req.RefreshToken)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while refreshing", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
