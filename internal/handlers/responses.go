package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/HeadGardener/medods/internal/services"
)

type response struct {
	Msg   string `json:"Msg"`
	Error string `json:"Error"`
}

func newErrResponse(w http.ResponseWriter, code int, msg string, err error) {
	log.Printf("[ERROR] %s: %s", msg, err.Error())
	if !errIsCustom(err) && code >= http.StatusInternalServerError {
		newResponse(w, code, response{
			Msg:   msg,
			Error: "unexpected error",
		})
		return
	}

	newResponse(w, code, response{
		Msg:   msg,
		Error: err.Error(),
	})
}

func newResponse(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func errIsCustom(err error) bool {
	if errors.Is(err, services.ErrNotSameRefreshToken) {
		return true
	}

	if errors.Is(err, services.ErrInvalidSession) {
		return true
	}

	if errors.Is(err, services.ErrRefreshTokenExpired) {
		return true
	}

	if errors.Is(err, ErrTokenNotString) {
		return true
	}

	return false
}
