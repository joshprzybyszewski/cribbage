package mysql

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct{}

func (s *interactionService) Get(id model.PlayerID) (interaction.PlayerMeans, error) {
	result := interaction.PlayerMeans{}

	return result, nil
}

func (s *interactionService) Create(pm interaction.PlayerMeans) error {
	return nil
}

func (s *interactionService) Update(pm interaction.PlayerMeans) error {
	return nil
}
