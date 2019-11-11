package play

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func cutPhase(g *model.Game) error {
	behindDealer := roundCutter(g)

	cutterAPI := pAPIs[behindDealer]
	cutterAPI.NotifyBlocking(model.CutCard, nil)

	return nil
}

func handleCut(g *model.Game, cutAction model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if cutAction.Overcomes != model.CutCard {
		return errors.New(`Does not attempt to cut deck`)
	}
	if err := isWaitingForPlayer(g, cutAction); err != nil {
		return err
	}
	if cutAction.ID != roundCutter(g) {
		return errors.New(`Wrong player is cutting`)
	}

	cda, ok := cutAction.Action.(model.CutDeckAction)
	if !ok {
		return errors.New(`tried dealing with a different action`)
	}

	if cda.Percentage < 0 || cda.Percentage > 1 {
		dAPI.NotifyBlocking(model.CutCard, `Needs cut value between 0 and 1`)
		return nil
	}

	if len(g.BlockingPlayers) != 1 {
		log.Errorf("Expected one blocker for cut, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, cutAction)

	// cut the deck
	return cut(g, cda.Percentage, pAPIs)
}

func cut(g *model.Game, cutPercent float64, pAPIs map[model.PlayerID]interaction.Player) error {
	c := g.Deck.CutDeck(cutPercent)

	if jack := model.NewCardFromString(`jh`); c.Value == jack.Value {
		// Check if the dealer was cut a jack
		addPoints(g, g.CurrentDealer, 2, pAPIs, `his nibs`)
	}

	g.CutCard = c

	for _, pAPI := range pAPIs {
		pAPI.NotifyMessage("Cut card", g.CutCard.String())
	}
}