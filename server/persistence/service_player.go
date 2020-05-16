package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type PlayerService interface {
	Get(id model.PlayerID) (model.Player, error)

	Create(p model.Player) error
	UpdateGameColor(id model.PlayerID, gID model.GameID, color model.PlayerColor) error

	BeginGame(gID model.GameID, ps []model.Player) error
}
