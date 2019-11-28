package play

import (
	"errors"
	"log"

	"github.com/joshprzybyszewski/cribbage/logic/pegging"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

var _ PhaseHandler = (*peggingHandler)(nil)

type peggingHandler struct{}

func (*peggingHandler) Start(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	g.PeggedCards = g.PeggedCards[:0]

	// put the player after the dealer as the blocking player
	pIDs := playersToDealTo(g)
	pID := pIDs[0]
	addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `please peg a card`)

	return nil
}

func (*peggingHandler) HandleAction(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if err := validateAction(g, action, model.PegCard); err != nil {
		return err
	}

	pa, ok := action.Action.(model.PegAction)
	if !ok {
		return errors.New(`tried pegging with a different action`)
	}

	pID := action.ID
	if pa.SayGo {
		curPeg := g.CurrentPeg()
		minCardVal := minUnpeggedValue(g.Hands[pID], g.PeggedCards)
		if curPeg+minCardVal <= model.MaxPeggingValue {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot say go when has unpegged playable card`)
			return nil
		}
	} else {
		if !handContains(g.Hands[pID], pa.Card) {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot peg card you don't have`)
			return nil
		}

		if hasBeenPegged(g.PeggedCards, pa.Card) {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot peg same card twice`)
			return nil
		}

		if g.CurrentPeg()+pa.Card.PegValue() > model.MaxPeggingValue {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot peg card with this value`)
			return nil
		}
	}

	if len(g.BlockingPlayers) != 1 {
		log.Printf("Expected one blocker for pegging, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, action)

	if pa.SayGo {
		err := doSayGo(g, action, pAPIs)
		if err != nil {
			return err
		}
	} else {
		err := doPeg(g, action, pa, pAPIs)
		if err != nil {
			return err
		}
	}

	if g.IsOver() {
		return nil
	}

	if len(g.PeggedCards) == 4*len(g.Players) {
		// This was the last card: give one point to this player.
		addPoints(g, action.ID, 1, pAPIs, `last card`)
		return nil
	}

	if g.IsOver() {
		return nil
	}

	// TODO this might need to be the next player with cards still to peg
	// Set the next player to peg as the blocker
	nextPlayerIndex := -1
	for i, p := range g.Players {
		if p.ID == action.ID {
			nextPlayerIndex = (i + 1) % len(g.Players)
			break
		}
	}

	bp := g.Players[nextPlayerIndex]
	addPlayerToBlocker(g, bp.ID, model.PegCard, pAPIs, ``)

	return nil
}

func doPeg(g *model.Game, action model.PlayerAction, pa model.PegAction, pAPIs map[model.PlayerID]interaction.Player) error {
	pts, err := pegging.PointsForCard(g.PeggedCards, pa.Card)
	if err != nil {
		return err
	}

	addPoints(g, action.ID, pts, pAPIs, `pegging`)

	g.PeggedCards = append(g.PeggedCards, model.PeggedCard{
		Card:     pa.Card,
		PlayerID: action.ID,
		Action:   g.NumActions() + 1,
	})

	return nil
}

func doSayGo(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if len(g.PeggedCards) == 0 {
		return nil
	}

	lastPeggerID := g.PeggedCards[len(g.PeggedCards)-1].PlayerID
	if lastPeggerID == action.ID {
		// The go's went all the way around. Take a point
		addPoints(g, action.ID, 1, pAPIs, `the go`)
	}

	return nil
}
