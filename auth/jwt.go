package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	username string
	jwt.StandardClaims
}

// Token is just used to serialize a token string to JSON
type Token struct {
	Token string `json:"token"`
}

// JWTService is used to generate tokens
type JWTService struct {
	jwtKey   []byte
	validFor time.Duration
}

func NewJWTService(key string, validFor time.Duration) JWTService {
	return JWTService{
		jwtKey:   []byte(key),
		validFor: validFor,
	}
}

func (js JWTService) CreateToken(username string) (Token, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(js.validFor).Unix(),
		},
	})
	tokStr, err := tok.SignedString(js.jwtKey)
	if err != nil {
		return Token{}, err
	}
	return Token{Token: tokStr}, nil
}
