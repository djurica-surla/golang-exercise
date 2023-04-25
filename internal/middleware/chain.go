package middleware

import (
	"net/http"
)

// Middleware type is used to describe the way middleware needs to be implemented
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Chain is a variadic function that chains and calls middlware from the first one provided is called first
func Chain(handler http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return handler
	}

	wrapped := handler

	for i := len(m) - 1; i >= 0; i-- {
		// middleware2(middleware1(handler))
		wrapped = m[i](wrapped)
	}

	return wrapped

}
