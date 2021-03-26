package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type InteractionService interface {
	Get(id model.PlayerID) (interaction.PlayerMeans, error)

	Create(pm interaction.PlayerMeans) error
	Update(pm interaction.PlayerMeans) error
}
