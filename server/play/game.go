package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func HandleAction(g *model.Game, action PlayerAction) (error) {
	if g.ID != action.GameID {
		return nil, errors.New(`action not for this game`)
	}

	canFulfill := false
	for _, bp := range g.BlockingPlayers {
		if bp.ID == action.ID && bp.Reason == action.Overcomes {
			canFulfill = true
		}
	}
	if !canFulfill {
		return nil, errors.New(`action does not overcome appropriate blocker`)
	}

	return g, nil
}

func PlayOneStep(g *model.Game) error {
	if g.IsOver() {
		return nil
	}

	playerAPIs := map[model.PlayerID]interaction.Player{}

	var err error
	switch g.Phase {
	case model.Deal:
		err = dealPhase(g, playerAPIs)
		g.Phase = model.BuildCrib
	case model.BuildCrib:
		err = buildCrib(g, playerAPIs)
		g.Phase = model.Cut
	case model.Cut:
		err = cutPhase(g, playerAPIs)
		g.Phase = model.Pegging
	case model.Pegging:
		err = peg(g, playerAPIs)
		g.Phase = model.Counting
	case model.Counting:
		err = countHands(g, playerAPIs)
		g.Phase = model.CribCounting
	case model.CribCounting:
		countCrib(g, playerAPIs)
		g.Phase = model.Done
		err = g.NextRound()
	}

	if err != nil {
		return err
	}

	return nil
}
