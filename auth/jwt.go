package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Username string `json:"username"`
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

// CreateToken creates and signs a JSON web token
func (js JWTService) CreateToken(username string) (Token, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Username: username,
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

// ValidateAndParseToken makes sure the token is valid given the token key, then returns the username contained within
// the token
func (js JWTService) ValidateAndParseToken(tokenStr string) (bool, string, error) {
	if tokenStr == `` {
		return false, ``, nil
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(`Unexpected signing method: %v`, token.Header[`alg`])
		}
		return js.jwtKey, nil
	})
	if err != nil {
		return false, ``, err
	}
	if cl, ok := tok.Claims.(*claims); ok && tok.Valid {
		return true, cl.Username, nil
	}
	return false, ``, nil
}
