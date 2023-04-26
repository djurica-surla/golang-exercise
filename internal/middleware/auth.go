package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/djurica-surla/golang-exercise/internal/token"
)

// TokenServicer represents necessary token service implementation for authentication middleware.
type TokenServicer interface {
	VerifyAccessToken(token string) error
}

// A middleware used for authorization of users.
func Authenticate(tokenService TokenServicer) func(next http.HandlerFunc) http.HandlerFunc {

	// Authenticate a user.
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the access cookie exists.
			accessToken, err := r.Cookie("accessToken")
			if err != nil {
				sendErrorReponse(w, "missing access token cookie, login again", http.StatusBadRequest)
				return
			}

			// Verify the access token, if there is an error the token is invalid.
			err = tokenService.VerifyAccessToken(accessToken.Value)
			if err != nil {
				// If its expired, tell the user to login again.
				if errors.Is(err, token.ErrExpiredToken) {
					sendErrorReponse(w, "login again! access token has expired", http.StatusBadRequest)
					return

				} else {
					// If access token is invalid, return error.
					sendErrorReponse(w, "could not verify access token cookie", http.StatusBadRequest)
					return
				}

			}
			next.ServeHTTP(w, r)
		})
	}
}

// sendErrorResponse encode the error and code into error response and sends it back.
func sendErrorReponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorResponse := ErrorResponse{
		Message: message,
		Code:    code,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
