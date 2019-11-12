package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type DB interface {
	GetGame(id model.GameID) (model.Game, error)
	GetPlayer(id model.PlayerID) (model.Player, error)
	GetInteraction(id model.PlayerID) (interaction.Player, error)

	SaveGame(g model.Game) error
	SavePlayer(p model.Player) error
	SaveInteraction(i interaction.Player) error
}
