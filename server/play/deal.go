package play

import (
	"errors"
	"log"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ PhaseHandler = (*dealingHandler)(nil)

type dealingHandler struct{}

func (*dealingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// Ensure all of the players hands are cleared before we start dealing
	for pID := range g.Hands {
		g.Hands[pID] = g.Hands[pID][:0]
	}

	// shuffle the deck at least once
	g.Deck.Shuffle()

	addPlayerToBlocker(g, g.CurrentDealer, model.DealCards, pAPIs, ``)

	return nil
}

func (*dealingHandler) HandleAction(g *model.Game,
	action model.PlayerAction,
	pAPIs map[model.PlayerID]interaction.Player,
) error {

	if err := validateAction(g, action, model.DealCards); err != nil {
		return err
	}

	if action.ID != g.CurrentDealer {
		return errors.New(`Wrong player is dealing`)
	}

	da, ok := action.Action.(model.DealAction)
	if !ok {
		return errors.New(`tried dealing with a different action`)
	}

	if da.NumShuffles <= 0 {
		addPlayerToBlocker(g, g.CurrentDealer, model.DealCards, pAPIs, `Need to shuffle at least once`)
		return nil
	}

	if len(g.BlockingPlayers) != 1 {
		log.Printf("Expected one blocker for deal, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, action)

	for i := 0; i < da.NumShuffles; i++ {
		// shuffle
		g.Deck.Shuffle()
	}

	// deal
	if err := deal(g, pAPIs); err != nil {
		return err
	}

	return nil
}

func deal(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// Get the order of players we need to deal to
	pIDs := playersToDealTo(g)

	// Define how many cards we need to deal and the hand size
	handSize := 6
	switch len(pIDs) {
	case 3, 4:
		handSize = 5
	}
	numCardsToDeal := handSize * len(pIDs)

	for numDealt := 0; numDealt < numCardsToDeal; {
		for _, pID := range pIDs {
			c := g.Deck.Deal()
			numDealt++
			g.Hands[pID] = append(g.Hands[pID], c)
		}
	}

	// For three player games, we need to deal another card to the crib
	if len(pIDs) == 3 {
		c := g.Deck.Deal()
		g.Crib = append(g.Crib, c)
	}

	// Now that the hands are all dealt, tell everyone about what they have
	for pID, hand := range g.Hands {
		handStr := handString(hand)
		_ = pAPIs[pID].NotifyMessage(*g, `Received Hand `+handStr)
	}

	return nil
}
