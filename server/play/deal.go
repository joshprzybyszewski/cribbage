package play

import (
	"log"

	"github.com/joshprzybyszewski/cribbage/model"
)

func dealPhase(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// shuffle the deck at least once
	g.Deck.Shuffle()

	dAPI := pAPIs[g.CurrentDealer]
	dAPI.NotifyBlocking(model.DealCards, nil)
	
	return nil
}

func handleDeal(g *model.Game, dealAction PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if dealAction.Overcomes != model.DealCards {
		return errors.New(`Does not attempt to deal`)
	}
	if err := isWaitingForPlayer(g, dealAction); err != nil {
		return err
	}
	if dealAction.ID != g.CurrentDealer {
		return errors.New(`Wrong player is dealing`)
	}

	da, ok := dealAction.(model.DealAction)
	if !ok {
		return errors.New(`tried dealing with a different action`)
	}

	if da.NumShuffles <= 0 {
		dAPI.NotifyBlocking(model.DealCards, `Need to shuffle at least once`)
		return nil
	}

	if len(g.BlockingPlayers) != 1 {
		log.Errorf("Expected one blocker for deal, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, dealAction)

	for i := 0; i < da.NumShuffles; i++ {
		// shuffle
		d.Shuffle()
	}

	// deal
	return deal(g, pAPIs)
}

func deal(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// Ensure all of the players hands are cleared before we start dealing
	for pID := range g.Hands {
		g.Hands[pID] = g.Hands[pID][:0]
	}
	g.Crib = g.Crib[:0]

	// Get the order of players we need to deal to
	pIDs := playersToDealTo(g)

	// Define how many cards we need to deal
	numCardsToDeal := 6 * 2
	if len(ps) == 3 {
		numCardsToDeal = 5 * 3
	} else if len(ps) == 4 {
		numCardsToDeal = 5 * 4
	}

	for numDealt := 0; numDealt < numCardsToDeal; {
		for _, pID := range pIDs {
			c, err := g.Deck.DealCard()
			if err != nil {
				return err
			}
			numDealt++
			g.Hands[pID] = append(g.Hands[pID], c)
		}
	}

	// For three player games, we need to deal another card to the crib
	if len(pIDs) == 3 {
		c, err := g.Deck.DealCard()
		if err != nil {
			return err
		}
		g.Crib = append(g.Crib, c)
	}

	// Now that the hands are all dealt, tell everyone about what they have
	for pID, hand := range g.Hands {
		handStr := ``
		for _, c := range hand {
			if len(handStr) > 0 {
				handStr += `, `
			}
			handStr += c.String()
		}
		pAPIs[pID].NotifyMessage(`Received Hand `+handStr)
	}

	return nil
}
