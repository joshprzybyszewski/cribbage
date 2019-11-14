package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type DB interface {
	CreatePlayer(p model.Player) error
	GetPlayer(id model.PlayerID) (model.Player, error)
	AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error

	GetGame(id model.GameID) (model.Game, error)
	SaveGame(g model.Game) error

	GetInteraction(id model.PlayerID) (interaction.Player, error)
	SaveInteraction(i interaction.Player) error
}
