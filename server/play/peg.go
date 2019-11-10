package play

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
)

func peg(g *Game) error {
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
