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
	// invalidate whatever card was cut last on this game
	g.CutCard = model.Card{}

	// set the blocker to be the player behind the dealer
	behindDealer := roundCutter(g)
	addPlayerToBlocker(g, behindDealer, model.CutCard, pAPIs, ``)

	return nil
}

func (*cuttingHandler) HandleAction(
	g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) error {

	cutPercent, err := cutActionValidation(g, action, pAPIs)
	if err != nil {
		return err
	}

	deck, err := g.GetDeck()
	if err != nil {
		return err
	}

	// cut the deck
	c, err := deck.CutDeck(cutPercent)
	if err != nil {
		return err
	}

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

func cutActionValidation(
	g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) (float64, error) {

	if err := validateAction(g, action, model.CutCard); err != nil {
		return 0, err
	}
	if action.ID != roundCutter(g) {
		return 0, errors.New(`Wrong player is cutting`)
	}

	cda, ok := action.Action.(model.CutDeckAction)
	if !ok {
		return 0, errors.New(`tried cutting with a different action`)
	}

	if cda.Percentage < 0 || cda.Percentage > 1 {
		addPlayerToBlocker(g, action.ID, model.CutCard, pAPIs, `Needs cut value between 0 and 1`)
		return 0, nil
	}

	if len(g.BlockingPlayers) != 1 {
		log.Printf("Expected one blocker for cut, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, action)

	return cda.Percentage, nil
}
