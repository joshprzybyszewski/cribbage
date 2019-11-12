package play

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type PhaseHandler interface {
	Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error
	HandleAction(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error
}

func validateAction(g *model.Game, action model.PlayerAction, blocker model.Blocker) error {
	if action.Overcomes != blocker {
		return fmt.Errorf(`Should overcome %v, but overcomes %v`, blocker, action.Overcomes)
	}
	if err := isWaitingForPlayer(g, action); err != nil {
		return err
	}

	return nil
}