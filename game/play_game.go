package game

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
)

func PlayEntireGame(g *Game) error {
	for !g.IsOver() {
		err := PlayOneStep(g)
		if err != nil {
			return err
		}
	}

	return nil
}

func PlayOneStep(g *Game) error {
	if g.IsOver() {
		return nil
	}

	var err error
	switch g.round.CurrentStage {
	case model.Deal:
		err = dealPhase(g)
		g.round.CurrentStage = model.BuildCrib
	case model.BuildCrib:
		err = buildCrib(g)
		g.round.CurrentStage = model.Cut
	case model.Cut:
		err = cutPhase(g)
		g.round.CurrentStage = model.Pegging
	case model.Pegging:
		err = peg(g)
		g.round.CurrentStage = model.Counting
	case model.Counting:
		err = countHands(g)
		g.round.CurrentStage = model.CribCounting
	case model.CribCounting:
		err = countCrib(g)
		if err != nil {
			return err
		}
		g.round.CurrentStage = model.DealingReady
		err = g.NextRound()
	}

	if err != nil {
		return err
	}

	return nil
}

func dealPhase(g *Game) error {
	// shuffle the deck at least once
	g.Deck().Shuffle()

	// init dealer
	d := g.Dealer()
	d.TakeDeck(g.Deck())

	// shuffle
	d.Shuffle()

	// deal
	return deal(g)
}

func deal(g *Game) error {
	d := g.Dealer()
	ps := g.PlayersToDealTo()

	numCardsToDeal := 6 * 2
	if len(ps) == 3 {
		numCardsToDeal = 5 * 3
	} else if len(ps) == 4 {
		numCardsToDeal = 4 * 4
	}
	for i := 0; i < numCardsToDeal; i++ {
		for _, p := range ps {
			c, err := d.DealCard()
			if err != nil {
				return err
			}
			err = p.AcceptCard(c)
			if err != nil {
				return err
			}
		}
	}

	// For three player games, we need to deal another card to the crib
	if len(ps) == 3 {
		c, err := d.DealCard()
		if err != nil {
			return err
		}
		err = g.round.AcceptCribCards(c)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func buildCrib(g *Game) error {
	ps := g.PlayersToDealTo()
	var wg sync.WaitGroup
	var err error

	for _, p := range ps {
		wg.Add(1)

		go func(pcopy Player) {
			desired := 2
			if len(ps) > 2 {
				desired = 1
			}
			err = g.round.AcceptCribCards(pcopy.AddToCrib(g.Dealer().Color(), desired)...)
			wg.Done()
		}(p)
	}

	wg.Wait()
	return err
}

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
