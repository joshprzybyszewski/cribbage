package game

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
)

type PegColor int

const (
	Blue PegColor = iota
	Red
	Green
)

const (
	winningScore int = 121
)

type GameConfig struct {
	Players        []Player
	StartingDealer int
	StartingCrib   []cards.Card
}

type Game struct {
	// The deck of cards
	deck *cards.Deck

	// The card that has been cut
	cutCard cards.Card
	hasCut  bool

	// The current "round" of cribbage
	round *Round

	// The dealer who also gets the crib
	dealer int

	// An ordered list of players
	players []Player

	// The current scores per color
	ScoresByColor map[PegColor]int

	// The previous scores per color
	LagScoreByColor map[PegColor]int
}

func New(cfg GameConfig) *Game {
	var r *Round
	switch len(cfg.Players) {
	case 2:
		r = NewTwoPlayerRound()
	default:
		return nil
	}

	return &Game{
		deck:            cards.NewDeck(),
		dealer:          cfg.StartingDealer,
		round:           r,
		players:         cfg.Players,
		ScoresByColor:   map[PegColor]int{Blue: 0, Red: 0},
		LagScoreByColor: map[PegColor]int{Blue: -1, Red: -1},
	}
}

func (g *Game) IsOver() bool {
	for _, score := range g.ScoresByColor {
		if score >= winningScore {
			return true
		}
	}
	return false
}

func (g *Game) Dealer() Player {
	return g.players[g.dealer]
}

func (g *Game) PlayersToDealTo() []Player {
	if g.dealer == len(g.players) - 1 {
		return g.players
	}

	return append(g.players[g.dealer+1:], g.players[:g.dealer+1]...)
}

func (g *Game) Deck() *cards.Deck {
	return g.deck
}

func (g *Game) CurrentRound() *Round {
	return g.round
}

func (g *Game) CutAt(p float64) error {
	if g.hasCut {
		return errors.New(`cannot re-cut the deck`)
	}

	g.cutCard = g.deck.CutDeck(p)

	g.hasCut = true
	return nil
}

func (g *Game) LeadCard() cards.Card {
	if !g.hasCut {
		return cards.Card{}
	}

	return g.cutCard
}

func (g *Game) NextRound() error {
	err := g.round.NextRound()
	if err != nil {
		return err
	}

	g.hasCut = false
	g.dealer = (g.dealer+1)%len(g.players)

	return nil
}

func (g *Game) AddPoints(pc PegColor, p int) {
	g.LagScoreByColor[pc] = g.ScoresByColor[pc]
	g.ScoresByColor[pc] = g.ScoresByColor[pc] + p
}
