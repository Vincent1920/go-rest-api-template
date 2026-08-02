package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret      = []byte("change-me")
	jwtExpireHours = 24
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func ConfigureJWT(secret string, expireHours int) {
	if strings.TrimSpace(secret) != "" {
		jwtSecret = []byte(secret)
	}
	if expireHours > 0 {
		jwtExpireHours = expireHours
	}
}

func GenerateToken(userID uint) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(jwtExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
