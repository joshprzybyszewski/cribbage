package game

import (
	"fmt"
	"sync"
)

func (g *Game) Play() error {
	deck := g.Deck()

	for !g.IsOver() {
		// shuffle the deck at least once
		deck.Shuffle()

		// init dealer
		d := g.Dealer()
		d.TakeDeck(deck)

		// shuffle
		d.Shuffle()

		// deal
		ps := g.PlayersToDealTo()
		err := deal(d, ps)
		if err != nil {
			return err
		}

		// start the round
		r := g.CurrentRound()

		// build crib
		g.buildCrib(r, ps)

		// cut
		r.CurrentStage = Cut
		behindDealer := ps[len(ps)-2]
		err = g.CutAt(behindDealer.Cut())
		if err != nil {
			return err
		}

		for _, p := range ps {
			p.TellAboutCut(g.LeadCard())
		}

		// peg
		err = g.peg(r, ps)
		if err != nil {
			return err
		}
		if g.IsOver() {
			break
		}

		// count hand
		fmt.Printf("Counting players hands\n")
		r.CurrentStage = Counting
		for _, p := range ps {
			msg, s := p.HandScore(g.LeadCard())
			g.AddPoints(p.Color(), s, msg)
			if g.IsOver() {
				break
			}
		}
		if g.IsOver() {
			break
		}

		// count crib
		fmt.Printf("Counting the crib\n")
		fmt.Printf("Crib is: %s %s %s %s\n", r.Crib()[0], r.Crib()[1], r.Crib()[2], r.Crib()[3])
		r.CurrentStage = CribCounting
		err = d.AcceptCrib(r.Crib())
		if err != nil {
			return err
		}
		msg, pts, err := d.CribScore(g.LeadCard())
		if err != nil {
			return err
		}
		g.AddPoints(d.Color(), pts, msg)
		if g.IsOver() {
			break
		}

		r.CurrentStage = Done

		// progress the round
		err = g.NextRound()
		if err != nil {
			return err
		}
	}

	return nil
}

func deal(d Player, ps []Player) error {
	for everyoneIsHappy := false; !everyoneIsHappy; {
		everyoneIsHappy = true
		for _, p := range ps {
			c, err := d.DealCard()
			if err != nil {
				return err
			}
			p.AcceptCard(c)
			if p.NeedsCard() {
				everyoneIsHappy = false
			}
		}
	}

	return nil
}

func (g *Game) buildCrib(r *Round, ps []Player) error {
	r.CurrentStage = BuildCrib
	var wg sync.WaitGroup
	var err error

	for _, p := range ps {
		wg.Add(1)

		go func(pcopy Player) {
			desired := 2
			if len(ps) > 2 {
				desired = 1
			}
			err = r.AcceptCribCards(pcopy.AddToCrib(g.Dealer().Color(), desired))
			wg.Done()
		}(p)
	}

	wg.Wait()
	return err
}

func (g *Game) peg(r *Round, ps []Player) error {
	r.CurrentStage = Pegging
	someoneCanPlay := true
	var lastPegger Player
	for someoneCanPlay {
		someoneCanPlay = false
		for _, p := range ps {
			c, sayGo, canPlay := p.Peg(r.PrevPeggedCards(), r.CurrentPeg())
			if !canPlay {
				if lastPegger == p && canPlay {
					// the goes went all the way around -- take a point
					r.GoAround()
					g.AddPoints(p.Color(), 1, `the go`)
					if g.IsOver() {
						return nil
					}
				}
				continue
			}
			someoneCanPlay = true
			if sayGo {
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
