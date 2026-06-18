package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager issues and verifies signed JWT access tokens.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenManager creates a TokenManager with the given signing secret and token lifetime.
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

// Issue returns a signed JWT whose subject is the user ID.
func (m *TokenManager) Issue(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secret)
}

// Verify validates the token signature and expiry and returns the user ID.
func (m *TokenManager) Verify(tokenString string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject: %w", err)
	}

	return id, nil
}
