package play

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/logic/pegging"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func startPeggingPhase(g *model.Game, pAPIs map[model.PlayerID]interaction.Player) error {
	// TODO before the pegging round starts, clear out the g.PeggedCards var

	return nil
}

func handlePeg(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if action.Overcomes != model.PegCard {
		return errors.New(`Does not attempt to peg`)
	}
	if err := isWaitingForPlayer(g, action); err != nil {
		return err
	}

	pa, ok := action.Action.(model.PegAction)
	if !ok {
		return errors.New(`tried dealing with a different action`)
	}

	pID := action.ID
	if pa.SayGo {
		curPeg := g.CurrentPeg()
		minCardVal := minUnpeggedValue(g.Hands[pID], g.PeggedCards)
		if curPeg + minCardVal <= model.MaxPeggingValue {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot say go when has unpegged playable card`)
			return nil
		}
	} else {
		if handContains(g.Hands[pID], pa.Card) {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot peg card you don't have`)
			return nil
		}
	
		if hasBeenPegged(g.PeggedCards, pa.Card) {
			addPlayerToBlocker(g, pID, model.PegCard, pAPIs, `Cannot peg same card twice`)
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

	if len(g.PeggedCards) == 4 * len(g.Players) {
		// This was the last card: give one point to this player.
		addPoints(g, action.ID, 1, pAPIs, `last card`)
		// TODO ensure the phase moves on to counting
		return nil
	}

	if g.IsOver() {
		return nil
	}

	// Set the next player to peg as the blocker
	nextPlayerIndex := -1
	for i, pID := range g.Players {
		if pID == action.ID {
			nextPlayerIndex = (i + 1) % len(g.Players)
			break
		}
	}

	bpID := g.Players[nextPlayerIndex]
	addPlayerToBlocker(g, bpID, model.PegCard,pAPIs)

	return nil
}

func doPeg(g *model.Game, action model.PlayerAction, pa model.PegAction, pAPIs map[model.PlayerID]interaction.Player) error {
	pIDs := playersToDealTo(g)

	pts, err := pegging.PointsForCard(g.PeggedCards, pa.Card)
	if err != nil {
		return 0, err
	}

	addPoints(g, action.ID, pts, pAPIs, `pegging`)

	g.PeggedCards = append(g.PeggedCards, model.PeggedCard{
		Card: pa.Card,
		PlayerID: action.ID,
	})

	return nil
}

func doSayGo(g *model.Game, action model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	if len(g.PeggedCards) == 0 {
		return nil
	}
	
	lastPeggerID := g.PeggedCards[len(g.PeggedCards)-1].ID
	if lastPeggerID == action.ID {
		// The go's went all the way around. Take a point
		addPoints(g, action.ID, 1, pAPIs, `the go`)
	}

	return nil
}