package play

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func peg(g *model.Game) error {
	r := g.round
	ps := g.PlayersToDealTo()
	var lastPegger Player
	for len(r.PrevPeggedCards()) < 4*len(ps) {
		for _, p := range ps {
			c, sayGo, canPlay := p.Peg(r.PrevPeggedCards(), r.CurrentPeg())
			if !canPlay || sayGo {
				if lastPegger == p {
					// the goes went all the way around -- take a point
					r.GoAround()
					g.AddPoints(p.Color(), 1, `the go`)
					if g.IsOver() {
						return nil
					}
				}
				continue
			}

			lastPegger = p
			pts, err := r.AcceptPegCard(c)
			if err != nil {
				return err
			}

			g.AddPoints(p.Color(), pts, `pegging`)
			if g.IsOver() {
				return nil
			}
		}
	}

	// give a point for last card
	g.AddPoints(lastPegger.Color(), 1, `last card`)
	if g.IsOver() {
		return nil
	}

	return nil
}

func handlePeg(g *model.Game, pegAction model.PlayerAction, pAPIs map[model.PlayerID]interaction.Player) error {
	// TODO before the pegging round starts, clear out the g.PeggedCards var
	if pegAction.Overcomes != model.PegCard {
		return errors.New(`Does not attempt to peg`)
	}
	if err := isWaitingForPlayer(g, pegAction); err != nil {
		return err
	}

	pa, ok := pegAction.Action.(model.PegAction)
	if !ok {
		return errors.New(`tried dealing with a different action`)
	}

	pID := pegAction.ID
	pAPI := pAPIs[pID]
	if pa.SayGo {
		// TODO check if the player has a card in their hand that they can play that hasn't been pegged
	} else {
		if handContains(g.Hands[pID], pa.Card) {
			pAPI.NotifyBlocking(model.PegCard, `Cannot peg card you don't have`)
			return nil
		}
	
		if hasBeenPegged(g.PeggedCards, pa.Card) {
			pAPI.NotifyBlocking(model.PegCard, `Cannot peg same card twice`)
			return nil
		}
	}
	

	if len(g.BlockingPlayers) != 1 {
		log.Printf("Expected one blocker for pegging, but had: %+v\n", g.BlockingPlayers)
	}
	removePlayerFromBlockers(g, pegAction)

	// TODO logic and add new blocker
}
