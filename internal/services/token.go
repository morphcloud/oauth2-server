package services

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token interface {
	Generate(tokenType string, expiresIn int64) (string, error)
}

type accessTokenCustomClaims struct {
	userID       uint64
	userRoleID   uint32
	userStatusID uint32
	jwt.StandardClaims
}

type refreshTokenCustomClaims struct {
	jwt.StandardClaims
}

type jwtToken struct {}

func NewJWTToken() *jwtToken {
	return &jwtToken{}
}

// newToken generates JWT token with custom claims and returns token string
func (t *jwtToken) Generate(tokenType string, expiresIn int64) (string, error) {
	if tokenType == "access_token" {
		accessTokenClaims := &accessTokenCustomClaims{
			userID:         512,
			userRoleID:     2,
			userStatusID:   1,
			StandardClaims: jwt.StandardClaims{
				Issuer: os.Getenv("APP_NAME"),
				IssuedAt: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + expiresIn,
			},
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
		accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return "", errors.New(err.Error())
		}

		return accessTokenString, nil
	} else if tokenType == "refresh_token" {
		refreshTokenClaims := &refreshTokenCustomClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer: os.Getenv("APP_NAME"),
				IssuedAt: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + expiresIn,
			},
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)
		refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return "", errors.New(err.Error())
		}

		return refreshTokenString, nil
	}

	return "", errors.New("invalid token type")
}
