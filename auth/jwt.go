package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	username string
	jwt.StandardClaims
}

type JWTService struct {
	jwtKey   string
	validFor time.Duration
}

func NewJWTService(key string, validFor time.Duration) JWTService {
	return JWTService{
		jwtKey:   key,
		validFor: validFor,
	}
}

func (js JWTService) CreateToken(username string) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(js.validFor).Unix(),
		},
	})
	tokStr, err := tok.SignedString(js.jwtKey)
	if err != nil {
		return ``, err
	}
	return tokStr, nil
}
