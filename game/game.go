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

type Game struct {
	// The deck of cards
	deck *cards.Deck

	// The current "round" of cribbage
	round *Round

	// The dealer who also gets the crib
	roundOwner *Player

	// An ordered list of players
	players []*Player

	// The current scores per color
	ScoresByColor map[PegColor]int

	// The previous scores per color
	LagScoreByColor map[PegColor]int
}

func New() *Game {
	human := NewPlayer(`human`)
	npc := NewPlayer(`computer`)

	return &Game{
		deck:            cards.NewDeck(),
		roundOwner:      human,
		round:           NewTwoPlayerRound(),
		players:         []*Player{human, npc},
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
