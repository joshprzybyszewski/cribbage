package game

import (
	"errors"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/scorer"
)

type Player interface {
	Name() string
	Color() PegColor
	TellAboutScores(cur, lag map[PegColor]int, msgs ...string)

	TakeDeck(d *cards.Deck)
	IsDealer() bool
	Shuffle()

	DealCard() (cards.Card, error)
	AcceptCard(cards.Card) error
	NeedsCard() bool

	AddToCrib(dealer PegColor, desired int) []cards.Card
	AcceptCrib([]cards.Card) error

	Cut() float64
	TellAboutCut(cards.Card)

	Peg(prevPegs []cards.Card, curPeg int) (played cards.Card, sayGo, canPlay bool)

	HandScore(leadCard cards.Card) (string, int)

	CribScore(leadCard cards.Card) (string, int, error)

	ReturnCards() error
}

type player struct {
	name        string
	color       PegColor
	interaction PlayerInteraction

	deck *cards.Deck

	hand   []cards.Card
	pegged map[cards.Card]struct{}
	crib   []cards.Card

	scoresByColor   map[PegColor]int
	lagScoreByColor map[PegColor]int
}

func newPlayer(interaction PlayerInteraction, name string, color PegColor) *player {
	return &player{
		interaction: interaction,
		name:        name,
		color:       color,
		deck:        nil,
		hand:        make([]cards.Card, 0, 6),
		crib:        make([]cards.Card, 0, 4),
		pegged:      make(map[cards.Card]struct{}, 4),
	}
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Color() PegColor {
	return p.color
}

func (p *player) TellAboutScores(cur, lag map[PegColor]int, msgs ...string) {
	p.scoresByColor = cur
	p.lagScoreByColor = lag
	p.interaction.TellAboutScores(cur, lag, msgs...)
}

func (p *player) TakeDeck(d *cards.Deck) {
	p.deck = d
}

func (p *player) IsDealer() bool {
	return p.deck != nil
}

func (p *player) DealCard() (cards.Card, error) {
	if !p.IsDealer() {
		return cards.Card{}, errors.New(`cannot deal a card if not the dealer`)
	}

	return p.deck.Deal(), nil
}

func (p *player) AcceptCard(c cards.Card) error {
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

func (p *player) AcceptCrib(crib []cards.Card) error {
	if !p.IsDealer() {
		return errors.New(`attempted to receive a crib when not the dealer`)
	}

	p.crib = crib

	return nil
}

func (p *player) HandScore(leadCard cards.Card) (string, int) {
	msg := fmt.Sprintf("hand (%s %s %s %s) with lead (%s)", p.hand[0], p.hand[1], p.hand[2], p.hand[3], leadCard)
	return msg, scorer.CribPoints(leadCard, p.hand)
}

func (p *player) CribScore(leadCard cards.Card) (string, int, error) {
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

func (p *player) AddToCrib(dealerColor PegColor, desired int) []cards.Card {
	cribCards := p.interaction.AskForCribCards(dealerColor, desired, p.hand)
	if len(cribCards) != desired {
		println(`bad time! never choose more than 2 cards`)
		return nil
	}

	// remove those cards from our hand
	p.hand = removeCards(p.hand, cribCards)

	return cribCards
}

func removeCards(before, without []cards.Card) []cards.Card {
	cc := without[0]
	after := make([]cards.Card, 0, len(before))
	for i, c := range before {
		if c.String() == cc.String() {
			if i == 0 {
				after = before[1:]
			} else if i == len(before)-1 {
				after = before[:i]
			} else {
				after = append(before[0:i], before[i+1:]...)
			}
		}
	}
	if len(without) == 1 {
		return after
	}

	cc = without[1]
	for i, c := range after {
		if c.String() == cc.String() {
			if i == 0 {
				after = after[1:]
			} else if i == len(after)-1 {
				after = after[:i]
			} else {
				after = append(after[0:i], after[i+1:]...)
			}
		}
	}
	return after
}

func (p *player) Cut() float64 {
	return p.interaction.AskForCut()
}

func (p *player) TellAboutCut(c cards.Card) {
	p.interaction.TellAboutCut(c)
}

func (p *player) Peg(prevPegs []cards.Card, curPeg int) (cards.Card, bool, bool) {
	if len(p.pegged) == len(p.hand) {
		return cards.Card{}, false, false
	}

	opts := make([]cards.Card, 0, len(p.hand))
	for _, c := range p.hand {
		if _, ok := p.pegged[c]; !ok {
			opts = append(opts, c)
		}
	}
	var c cards.Card
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
		return cards.Card{}, true, true
	}
	p.pegged[c] = struct{}{}

	return c, false, true
}
