package play

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/game"
)

func PlayGame() error {
	human := game.NewHumanPlayer(game.Blue)
	npc := game.NewDumbNPC(game.Red)
	cfg := game.GameConfig{
		Players:        []game.Player{human, npc},
		StartingDealer: human,
	}

	g := game.New(cfg)
	deck := g.Deck()

	for !g.IsOver() {
		// init dealer
		d := g.Dealer()
		d.TakeDeck(deck)

		// shuffle
		println(`shuffling`)
		d.Shuffle()

		// deal
		println(`dealing`)
		ps := g.PlayersToDealTo()
		err := deal(d, ps)
		if err != nil {
			return err
		}

		// start the round
		r := g.CurrentRound()

		// build crib
		println(`building crib`)
		buildCrib(g, r, ps)

		// cut
		println(`cutting`)
		behindDealer := ps[len(ps)-2]
		err = g.CutAt(behindDealer.Cut())
		if err != nil {
			return err
		}

		// peg
		println(`pegging`)
		peg(g, r, ps)

		// count
		for _, p := range ps {
			s := p.HandScore(g.LeadCard())
			g.AddPoints(p.Color(), s)
			// TODO check for a winner
		}

		// count crib
		println(`counting crib`)
		d.AcceptCrib(r.Crib())
		d.CribScore(g.LeadCard())

		// progress the round
		println(`progressing`)
		err = g.NextRound()
		if err != nil {
			return err
		}
	}

	return nil
}

func deal(d game.Player, ps []game.Player) error {
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

func buildCrib(g *game.Game, r *game.Round, ps []game.Player) error {
	var wg sync.WaitGroup
	var err error

	for _, p := range ps {
		wg.Add(1)

		go func(pcopy game.Player) {
			desired := 2
			if len(ps) > 2 {
				desired = 1
			}
			err = r.AcceptCribCards(pcopy.AddToCrib(g.Dealer().Color(), desired))
			println(`completed for ` + pcopy.Name())
			wg.Done()
		}(p)
	}

	wg.Wait()
	return err
}

func peg(g *game.Game, r *game.Round, ps []game.Player) error {
	someoneCanPlay := true
	var lastPegger game.Player
	for someoneCanPlay {
		someoneCanPlay = false
		for _, p := range ps {
			c, sayGo, canPlay := p.Peg(r.PrevPeggedCards(), r.CurrentPeg())
			if canPlay {
				someoneCanPlay = true
			} else {
				continue
			}
			if sayGo {
				if lastPegger == p {
					// the goes went all the way around -- take a point
					g.AddPoints(p.Color(), 1)
				}
				continue
			}

			lastPegger = p
			pts, err := r.AcceptPegCard(c)
			if err != nil {
				return err
			}

			g.AddPoints(p.Color(), pts)
		}
	}

	// give a point for last card
	g.AddPoints(lastPegger.Color(), 1)

	return nil
}
