package middleware

import (
	"encoding/json"
	"net/http"

	httpErrors "github.com/djurica-surla/golang-exercise/internal/errors"
)

// Error response represents how errors will be returned.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Custom handler describes how handler func needs to be implemented to be handled with HandleError.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// eandleError wraps around the custom handler func and resolves error type to a proper response.
func HandleError(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := f(w, r)
		if err != nil {
			var errorResponse ErrorResponse

			if httpErrors.IsUnauthorizedError(err) {
				w.WriteHeader(http.StatusUnauthorized)
				errorResponse.Code = http.StatusUnauthorized
			} else if httpErrors.IsBadRequestError(err) {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse.Code = http.StatusBadRequest
			} else if httpErrors.IsNotFoundError(err) {
				w.WriteHeader(http.StatusNotFound)
				errorResponse.Code = http.StatusNotFound
			} else if httpErrors.IsForbiddenError(err) {
				w.WriteHeader(http.StatusForbidden)
				errorResponse.Code = http.StatusForbidden
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				errorResponse.Code = http.StatusInternalServerError
			}

			errorResponse.Message = err.Error()
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
