package game

import (
	"errors"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

type Player interface {
	Name() string
	Color() model.PlayerColor
	TellAboutScores(cur, lag map[model.PlayerColor]int, msgs ...string)

	TakeDeck(d model.Deck)
	IsDealer() bool
	Shuffle()

	DealCard() (model.Card, error)
	AcceptCard(model.Card) error
	NeedsCard() bool

	AddToCrib(dealer model.PlayerColor, desired int) []model.Card
	AcceptCrib([]model.Card) error

	Cut() float64
	TellAboutCut(model.Card)

	Peg(prevPegs []model.PeggedCard, curPeg int) (played model.Card, sayGo, canPlay bool)

	HandScore(leadCard model.Card) (string, int)

	CribScore(leadCard model.Card) (string, int, error)

	ReturnCards() error
}

func NewPlayerFromModel(mp model.Player) Player {
	var interaction PlayerInteraction // TODO figure this out
	return &player{
		interaction: interaction,
		name:        mp.Name,
		color:       mp.Color,
		deck:        nil,
		hand:        make([]model.Card, 0, 6),
		crib:        make([]model.Card, 0, 4),
		pegged:      make(map[model.Card]struct{}, 4),
	}
}

type player struct {
	name        string
	color       model.PlayerColor
	interaction PlayerInteraction

	deck model.Deck

	hand   []model.Card
	pegged map[model.Card]struct{}
	crib   []model.Card

	scoresByColor   map[model.PlayerColor]int
	lagScoreByColor map[model.PlayerColor]int
}

func newPlayer(interaction PlayerInteraction, name string, color model.PlayerColor) *player {
	return &player{
		interaction: interaction,
		name:        name,
		color:       color,
		deck:        nil,
		hand:        make([]model.Card, 0, 6),
		crib:        make([]model.Card, 0, 4),
		pegged:      make(map[model.Card]struct{}, 4),
	}
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Color() model.PlayerColor {
	return p.color
}

func (p *player) TellAboutScores(cur, lag map[model.PlayerColor]int, msgs ...string) {
	p.scoresByColor = cur
	p.lagScoreByColor = lag
	p.interaction.TellAboutScores(cur, lag, msgs...)
}

func (p *player) TakeDeck(d model.Deck) {
	p.deck = d
}

func (p *player) IsDealer() bool {
	return p.deck != nil
}

func (p *player) DealCard() (model.Card, error) {
	if !p.IsDealer() {
		return model.Card{}, errors.New(`cannot deal a card if not the dealer`)
	}

	return p.deck.Deal(), nil
}

func (p *player) AcceptCard(c model.Card) error {
	if !p.NeedsCard() {
		return errors.New(`cannot accept new card`)
	}

	p.hand = append(p.hand, c)
	return nil
}

func (p *player) ReturnCards() error {
	if len(p.hand) == 0 || len(p.pegged) == 0 {
		return errors.New(`no cards to relinquish`)
	}
	p.hand = p.hand[:0]
	p.crib = p.crib[:0]
	for k := range p.pegged {
		delete(p.pegged, k)
	}

	return nil
}

func (p *player) NeedsCard() bool {
	// TODO this is only currently setup for 2 player games
	return len(p.hand) < 6
}

func (p *player) AcceptCrib(crib []model.Card) error {
	if !p.IsDealer() {
		return errors.New(`attempted to receive a crib when not the dealer`)
	}

	p.crib = crib

	return nil
}

func (p *player) HandScore(leadCard model.Card) (string, int) {
	msg := fmt.Sprintf("hand (%s %s %s %s) with lead (%s)", p.hand[0], p.hand[1], p.hand[2], p.hand[3], leadCard)
	return msg, scorer.HandPoints(leadCard, p.hand)
}

func (p *player) CribScore(leadCard model.Card) (string, int, error) {
	if !p.IsDealer() {
		return ``, 0, errors.New(`Cannot score crib when not the dealer`)
	} else if len(p.crib) == 0 {
		return ``, 0, errors.New(`expected to have crib, but missing!`)
	}

	msg := fmt.Sprintf("crib (%s %s %s %s) with lead (%s)", p.crib[0], p.crib[1], p.crib[2], p.crib[3], leadCard)
	return msg, scorer.CribPoints(leadCard, p.crib), nil
}

// interactions

func (p *player) Shuffle() {
	for p.interaction.AskToShuffle() {
		p.deck.Shuffle()
	}
}

func (p *player) AddToCrib(dealerColor model.PlayerColor, desired int) []model.Card {
	cribCards := p.interaction.AskForCribCards(dealerColor, desired, p.hand)
	if len(cribCards) != desired {
		fmt.Printf(`bad time! Expected %d cards chosen, but was %d (%v)\n`, desired, len(cribCards), cribCards)
		return nil
	}

	inCrib := map[model.Card]struct{}{}
	for _, cc := range cribCards {
		inCrib[cc] = struct{}{}
	}
	shouldRemove := func(c model.Card) bool {
		_, ok := inCrib[c]
		return ok
	}

	// remove those cards from our hand
	p.hand = removeCards(p.hand, shouldRemove)

	return cribCards
}

func removeCards(before []model.Card, shouldRemove func(model.Card) bool) []model.Card {
	after := make([]model.Card, 0, len(before))
	for _, c := range before {
		if !shouldRemove(c) {
			after = append(after, c)
		}
	}
	return after
}

func (p *player) Cut() float64 {
	return p.interaction.AskForCut()
}

func (p *player) TellAboutCut(c model.Card) {
	p.interaction.TellAboutCut(c)
}

func (p *player) Peg(prevPegs []model.PeggedCard, curPeg int) (model.Card, bool, bool) {
	if len(p.pegged) == len(p.hand) {
		return model.Card{}, false, false
	}

	opts := make([]model.Card, 0, len(p.hand))
	for _, c := range p.hand {
		if _, ok := p.pegged[c]; !ok {
			opts = append(opts, c)
		}
	}
	var c model.Card
	var sayGo bool
	for i := 0; ; i++ {
		if i == 10 {
			panic(errors.New(`the user is being very dumb`))
		}
		c, sayGo = p.interaction.AskToPeg(opts, prevPegs, curPeg)
		if sayGo {
			for _, o := range opts {
				if o.PegValue() <= maxPeggingValue-curPeg {
					continue
				}
			}
		} else {
			if c.PegValue() > maxPeggingValue-curPeg {
				continue
			}
		}

		break
	}
	if sayGo {
		return model.Card{}, true, true
	}
	p.pegged[c] = struct{}{}

	return c, false, true
}
