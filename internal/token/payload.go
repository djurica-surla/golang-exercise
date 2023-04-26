package token

import (
	"errors"
	"time"
)

var ErrExpiredToken = errors.New("token has expired")

// AccessTokenPayload contains the payload data of the access token.
type AccessTokenPayload struct {
	TokenID   string
	Username  string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

// Valid checks if token payload is valid or not.
func (payload *AccessTokenPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
