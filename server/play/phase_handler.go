package play

import (
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type PhaseHandler interface {
	// Start will do the one-time set up for this phase, alerting any players
	// if they are blocking, and assumes you will increment the phase after this call
	Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error
	// HandleAction will perform any validation on the action and/or game state, alerting any
	// blocking players, choosing to not error in favor of re-alerting the current blocker of
	// a misaction
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
