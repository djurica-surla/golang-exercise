package middleware

import (
	"encoding/json"
	"net/http"
)

// Custom handler describes how handler func needs to be implemented to be handled with HandleError.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// eandleError wraps around the custom handler func and resolves error type to a proper response.
func HandleError(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := f(w, r)
		if err != nil {
			if IsUnauthorizedError(err) {
				w.WriteHeader(http.StatusUnauthorized)
			} else if IsBadRequestError(err) {
				w.WriteHeader(http.StatusBadRequest)
			} else if IsNotFoundError(err) {
				w.WriteHeader(http.StatusNotFound)
			} else if IsForbiddenError(err) {
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			errorResponse := map[string]string{"error": err.Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
