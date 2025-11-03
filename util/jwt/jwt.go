package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignHS256(secret string, sub uint, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
