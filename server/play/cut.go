package play

import (
	"errors"
	"log"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ PhaseHandler = (*cuttingHandler)(nil)

type cuttingHandler struct{}

func (*cuttingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {

	behindDealer := roundCutter(g)

	addPlayerToBlocker(g, behindDealer, model.CutCard, pAPIs)

	return nil
}

func (*cuttingHandler) HandleAction(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if err := validateAction(g, action, model.CutCard); err != nil {
		return err
	}
	if action.ID != roundCutter(g) {
		return errors.New(`Wrong player is cutting`)
	}

	cda, ok := action.Action.(model.CutDeckAction)
	if !ok {
		return errors.New(`tried dealing with a different action`)
	}

	if cda.Percentage < 0 || cda.Percentage > 1 {
		addPlayerToBlocker(g, action.ID, model.CutCard, pAPIs, `Needs cut value between 0 and 1`)
		return nil
	}

	if len(g.BlockingPlayers) != 1 {
		log.Printf("Expected one blocker for cut, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, action)

	// cut the deck
	cutPercent := cda.Percentage
	c := g.Deck.CutDeck(cutPercent)

	if c.Value == model.JackValue {
		// Check if the dealer was cut a jack
		addPoints(g, g.CurrentDealer, 2, pAPIs, `his nibs`)
	}

	g.CutCard = c

	for _, pAPI := range pAPIs {
		_ = pAPI.NotifyMessage(*g, "Cut card "+g.CutCard.String())
	}

	return nil
}
