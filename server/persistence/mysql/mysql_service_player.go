package mysql

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct{}

func (ps *playerService) Create(p model.Player) error {
	// check if the player already exists
	return nil
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	return model.Player{}, nil
}

func (ps *playerService) UpdateGameColor(id model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	return nil
}
