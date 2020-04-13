package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Credentials are used to authorize a user
type Credentials struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

// NewCredentials takes a username and plaintext password, hashes the password, and returns Credentials
func NewCredentials(username, password string) (Credentials, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return Credentials{}, err
	}
	return Credentials{
		Username:     username,
		PasswordHash: string(hash),
	}, nil
}

// ValidateCredentials compares a username and plaintext password to a Credentials object
func ValidateCredentials(username, password string, creds Credentials) (bool, error) {
	if username != creds.Username {
		return false, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
