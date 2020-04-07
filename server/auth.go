package server

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

type AuthService struct {
	db persistence.DB
}

func NewAuthService(db persistence.DB) AuthService {
	return AuthService{
		db: db,
	}
}

func (as *AuthService) RegisterUser(username, password string) error {
	if len(password) < 6 {
		return errors.New(`password must be at least 6 characters`)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	player := model.Player{
		Auth: model.Auth{
			Username: username,
			Password: string(hash),
		},
		Name: username,
	}
	return as.db.CreatePlayer(player)
}

func (as *AuthService) LoginUser(username, password string) error {
	// TODO is the PlayerID the username?
	player, err := as.db.GetPlayer(model.PlayerID(username))
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(player.Auth.Password), []byte(password))
}
