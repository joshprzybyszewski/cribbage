package play

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
)

func countHands(g *Game) error {
	ps := g.PlayersToDealTo()
	for _, p := range ps {
		msg, s := p.HandScore(g.LeadCard())
		g.AddPoints(p.Color(), s, msg)
		if g.IsOver() {
			return nil
		}
	}

	return nil
}
