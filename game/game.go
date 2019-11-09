package game

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	winningScore int = 121
)

type GameConfig struct {
	Players        []Player
	StartingDealer int
	StartingCrib   []model.Card
}

type Game struct {
	// The deck of cards
	deck model.Deck

	// The card that has been cut
	cutCard model.Card
	hasCut  bool

	// The current "round" of cribbage
	round *Round

	// The dealer who also gets the crib
	dealer int

	// An ordered list of players
	players []Player

	// The current scores per color
	ScoresByColor map[model.PlayerColor]int

	// The previous scores per color
	LagScoreByColor map[model.PlayerColor]int
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
		deck:            model.NewDeck(),
		dealer:          cfg.StartingDealer,
		round:           r,
		players:         cfg.Players,
		ScoresByColor:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScoreByColor: map[model.PlayerColor]int{model.Blue: -1, model.Red: -1},
	}
}

func NewGameFromModel(mg model.Game) *Game {
	var r *Round
	switch len(mg.Players) {
	case 2:
		r = NewTwoPlayerRound()
	default:
		return nil
	}

	dealerIndex := 0
	gp := make([]Player, len(mg.Players))
	for i, p := range mg.Players {
		if p.ID == mg.CurrentDealer {
			dealerIndex = i
		}
		gp[i] = NewPlayerFromModel(p)
	}

	sbc := make(map[model.PlayerColor]int, len(mg.CurrentScores))
	for c, s := range mg.CurrentScores {
		sbc[c] = int(s)
	}

	lsbc := make(map[model.PlayerColor]int, len(mg.LagScores))
	for c, s := range mg.LagScores {
		lsbc[c] = int(s)
	}

	return &Game{
		deck:            model.NewDeck(),
		dealer:          dealerIndex,
		round:           r,
		players:         gp,
		ScoresByColor:   sbc,
		LagScoreByColor: lsbc,
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
	if g.dealer == len(g.players)-1 {
		return g.players
	}

	return append(g.players[g.dealer+1:], g.players[:g.dealer+1]...)
}

func (g *Game) Deck() model.Deck {
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

func (g *Game) LeadCard() model.Card {
	if !g.hasCut {
		return model.Card{}
	}

	return g.cutCard
}

func (g *Game) NextRound() error {
	err := g.round.NextRound()
	if err != nil {
		return err
	}

	g.hasCut = false
	g.dealer = (g.dealer + 1) % len(g.players)

	for _, p := range g.players {
		err = p.ReturnCards()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) AddPoints(pc model.PlayerColor, p int, msgs ...string) {
	g.LagScoreByColor[pc] = g.ScoresByColor[pc]
	g.ScoresByColor[pc] = g.ScoresByColor[pc] + p
	for _, p := range g.players {
		p.TellAboutScores(g.ScoresByColor, g.LagScoreByColor, msgs...)
	}
}
