package play

import (
	"fmt"
	"sync"

	"github.com/joshprzybyszewski/cribbage/game"
)

func PlayGame() error {
	human := game.NewHumanPlayer(game.Blue)
	npc := game.NewDumbNPC(game.Red)
	cfg := game.GameConfig{
		Players:        []game.Player{human, npc},
		StartingDealer: 0,
	}

	g := game.New(cfg)
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
		buildCrib(g, r, ps)

		// cut
		r.CurrentStage = game.Cut
		behindDealer := ps[len(ps)-2]
		err = g.CutAt(behindDealer.Cut())
		if err != nil {
			return err
		}

		for _, p := range ps {
			p.TellAboutCut(g.LeadCard())
		}

		// peg
		err = peg(g, r, ps)
		if err != nil {
			return err
		}
		if g.IsOver() {
			break
		}

		// count
		r.CurrentStage = game.Counting
		for _, p := range ps {
			s := p.HandScore(g.LeadCard())
			over := g.AddPoints(p.Color(), s)
			if over {
				break
			}
		}
		if g.IsOver() {
			break
		}

		// count crib
		r.CurrentStage = game.CribCounting
		err = d.AcceptCrib(r.Crib())
		if err != nil {
			return err
		}
		pts, err := d.CribScore(g.LeadCard())
		if err != nil {
			return err
		}
		over := g.AddPoints(d.Color(), pts)
		if over {
			break
		}

		r.CurrentStage = game.Done

		// progress the round
		err = g.NextRound()
		if err != nil {
			return err
		}
	}

	fmt.Printf("game over!\n")

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
	r.CurrentStage = game.BuildCrib
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
			wg.Done()
		}(p)
	}

	wg.Wait()
	return err
}

func peg(g *game.Game, r *game.Round, ps []game.Player) error {
	r.CurrentStage = game.Pegging
	someoneCanPlay := true
	var lastPegger game.Player
	for someoneCanPlay {
		someoneCanPlay = false
		for _, p := range ps {
			c, sayGo, canPlay := p.Peg(r.PrevPeggedCards(), r.CurrentPeg())
			if !canPlay {
				if lastPegger == p {
					// the goes went all the way around -- take a point
					r.GoAround()
					over := g.AddPoints(p.Color(), 1)
					if over {
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
					over := g.AddPoints(p.Color(), 1)
					if over {
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

			over := g.AddPoints(p.Color(), pts)
			if over {
				return nil
			}
		}
	}

	// give a point for last card
	over := g.AddPoints(lastPegger.Color(), 1)
	if over {
		return nil
	}

	return nil
}
