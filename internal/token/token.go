package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/djurica-surla/golang-exercise/internal/model"
)

const minSecretKeySize = 32

var ErrInvalidKeySize = fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
var ErrInvalidToken = errors.New("token is invalid")

// TokenService will create and verify jwt tokens.
type tokenService struct {
	secretKey string
}

// NewTokenMaker creates a new TokenMaker
func NewTokenService(secretKey string) (*tokenService, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInvalidKeySize
	}

	return &tokenService{
		secretKey: secretKey,
	}, nil
}

// CreateAccessToken creates an access token with the provided data
func (t *tokenService) CreateAccessToken(loginInfo model.Login) (string, error) {

	tokenID := uuid.New().String()

	payload := &AccessTokenPayload{
		TokenID:   tokenID,
		Username:  loginInfo.Username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute), // can be on config
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := jwtToken.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create access token: error creating token string: %w", err)
	}

	return tokenString, nil
}

// CreateAccessToken creates an access token with the provided data
func (t *tokenService) VerifyAccessToken(token string) error {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(t.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &AccessTokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			// If token is expired, still send the data back
			_, ok := jwtToken.Claims.(*AccessTokenPayload)
			if !ok {
				return ErrInvalidToken
			}
			return ErrExpiredToken
		}
		return ErrInvalidToken
	}

	// Convert the token to payload.
	// Payload can be returned and used, but in this case we just check if its valid.
	_, ok := jwtToken.Claims.(*AccessTokenPayload)
	if !ok {
		return ErrInvalidToken
	}

	return nil
}
