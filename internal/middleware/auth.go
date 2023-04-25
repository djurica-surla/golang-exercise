package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/djurica-surla/golang-exercise/internal/token"
)

// TokenServicer represents necessary token service implementation for authentication middleware.
type TokenServicer interface {
	VerifyAccessToken(token string) error
}

// A middleware used for authorization of users.
func Authenticate(tokenService TokenServicer) func(next http.HandlerFunc) http.HandlerFunc {

	// Authenticate a user
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Check if the access cookie exists
			accessToken, err := r.Cookie("accessToken")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse := map[string]string{"error": "missing access token cookie, login again"}
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			// Verify the access token, if there is an error the token is invalid.
			err = tokenService.VerifyAccessToken(accessToken.Value)
			if err != nil {
				// If its expired, validate refresh token
				if errors.Is(err, token.ErrExpiredToken) {
					w.WriteHeader(http.StatusBadRequest)
					errorResponse := map[string]string{"error": "login again! access token has expired"}
					json.NewEncoder(w).Encode(errorResponse)
					return

				} else {
					// If access token is invalid, return error
					w.WriteHeader(http.StatusUnauthorized)
					errorResponse := map[string]string{"error": fmt.Sprintf("could not verify access token cookie: %s", err.Error())}
					json.NewEncoder(w).Encode(errorResponse)
					return
				}

			}
			next.ServeHTTP(w, r)
		})
	}
}
