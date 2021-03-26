package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type GameService interface {
	Get(id model.GameID) (model.Game, error)
	GetAt(id model.GameID, numActions uint) (model.Game, error)

	UpdatePlayerColor(id model.GameID, pID model.PlayerID, color model.PlayerColor) error
	Begin(g model.Game) error
	Save(g model.Game) error
}
