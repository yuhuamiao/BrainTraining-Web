package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const developmentJWTSecret = "brain-training-local-development-secret"

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID string) (string, error) {
	now := time.Now()
	expireTime := now.Add(24 * time.Hour) // 有效期24小时

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "brain-training",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret())
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// ValidateJWTConfiguration prevents a production process from using a known or weak signing key.
func ValidateJWTConfiguration() error {
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "release") {
		return nil
	}
	configured := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if configured == "" || configured == developmentJWTSecret {
		return errors.New("JWT_SECRET must be set in release mode")
	}
	if len(configured) < 32 {
		return fmt.Errorf("JWT_SECRET must contain at least 32 characters in release mode")
	}
	return nil
}

func secret() []byte {
	if configured := strings.TrimSpace(os.Getenv("JWT_SECRET")); configured != "" {
		return []byte(configured)
	}
	return []byte(developmentJWTSecret)
}
