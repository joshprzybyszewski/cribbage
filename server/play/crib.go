package play

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func buildCrib(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// Clear out the previous crib before we start building this one
	g.Crib = g.Crib[:0]

	// Dell all of the players they need to give us the desired number of cards
	pIDs := playersToDealTo(g)
	desired := numDesiredCribCards(g)

	for _, pID := range pIDs {
		pAPIs[pID].NotifyBlocking(model.CribCard, model.CribBlocker{
			Desired: desired,
			Dealer: g.CurrentDealer,
			PlayerColors: g.PlayerColors,
		})
	}

	return nil
}

func handleCribBuild(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if action.Overcomes != model.CribCard {
		return errors.New(`Does not attempt to build crib`)
	}
	if err := isWaitingForPlayer(g, action); err != nil {
		return err
	}

	bca, ok := action.Action.(model.BuildCribAction)
	if !ok {
		return errors.New(`tried building crib with a different action`)
	}

	if len(bca.Cards) != numDesiredCribCards(g) {
		pAPIs[action.ID].NotifyBlocking(model.CribCard, `Need to submit all required cards at once`)
		return nil
	}
	if !isSuperSet(g.Hands[action.ID], bca.Cards) {
		pAPIs[action.ID].NotifyBlocking(model.CribCard, `Cannot submit cards that are not in your hand`)
		return nil
	}

	removePlayerFromBlockers(g, action)

	// Put the player's cards from their hand into the crib
	g.Crib = append(g.Crib, bca.Cards...)
	g.Hands[action.ID] = removeSubset(g.Hands[action.ID], bca.Cards)

	return nil
}

func numDesiredCribCards(g *model.Game) int {
	if len(g.Players) > 2 {
		return 1
	}
	return 2
}