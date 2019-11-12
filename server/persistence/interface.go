package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type Persistence interface {
	GetGame(id model.GameID) (model.Game, error)
	GetPlayer(id model.PlayerID) (model.Player, error)
	GetInteraction(id model.PlayerID) (interaction.Player, error)
}
