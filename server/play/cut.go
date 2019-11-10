package play

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
)

func cutPhase(g *Game) error {
	ps := g.PlayersToDealTo()
	behindDealer := ps[len(ps)-2]
	err := g.CutAt(behindDealer.Cut())
	if err != nil {
		return err
	}
	if g.LeadCard().Value == 11 {
		g.AddPoints(g.Dealer().Color(), 2, `nobs`)
	}

	for _, p := range ps {
		p.TellAboutCut(g.LeadCard())
	}

	return nil
}
