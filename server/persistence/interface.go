package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type Persistence interface {
	GetGame(id model.GameID) model.Game
	GetPlayer(id model.PlayerID) model.Game
}
