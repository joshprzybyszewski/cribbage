package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ PhaseHandler = (*cribBuildingHandler)(nil)
type cribBuildingHandler struct {}

func (*cribBuildingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// Clear out the previous crib before we start building this one
	g.Crib = g.Crib[:0]

	// Tell all of the players they need to give us the desired number of cards
	pIDs := playersToDealTo(g)
	desired := numDesiredCribCards(g)

	for _, pID := range pIDs {
		addPlayerToBlocker(g, pID, model.CribCard, pAPIs, model.CribBlocker{
			Desired: desired,
			Dealer: g.CurrentDealer,
			PlayerColors: g.PlayerColors,
		})
	}

	return nil
}

func (*cribBuildingHandler) HandleAction(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if err := validateAction(g, action, model.CribCard); err != nil {
		return err
	}

	bca, ok := action.Action.(model.BuildCribAction)
	if !ok {
		return errors.New(`tried building crib with a different action`)
	}

	if len(bca.Cards) != numDesiredCribCards(g) {
		addPlayerToBlocker(g, action.ID, model.CribCard, pAPIs,  `Need to submit all required cards at once`)
		return nil
	}
	if !isSuperSet(g.Hands[action.ID], bca.Cards) {
		addPlayerToBlocker(g, action.ID, model.CribCard, pAPIs, `Cannot submit cards that are not in your hand`)
		return nil
	}

	removePlayerFromBlockers(g, action)

	// Put the player's cards from their hand into the crib
	g.Crib = append(g.Crib, bca.Cards...)
	g.Hands[action.ID] = removeSubset(g.Hands[action.ID], bca.Cards)

	if len(g.BlockingPlayers) == 0 {
		if len(g.Crib) != 4 {
			return errors.New(`no remaining blockers, but not enough cards in the crib`)
		}
		g.Phase++
	}

	return nil
}

func numDesiredCribCards(g *model.Game) int {
	if len(g.Players) > 2 {
		return 1
	}
	return 2
}