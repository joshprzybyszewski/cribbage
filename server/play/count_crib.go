package play

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
)

func countCrib(g *Game) error {
	r := g.round
	d := g.Dealer()
	err := d.AcceptCrib(r.Crib())
	if err != nil {
		return err
	}
	msg, pts, err := d.CribScore(g.LeadCard())
	if err != nil {
		return err
	}
	g.AddPoints(d.Color(), pts, msg)

	return nil
}
