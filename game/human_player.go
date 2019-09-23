package game

import (
	"fmt"

	"github.com/AlecAivazis/survey"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var _ Player = (*humanPlayer)(nil)

type humanPlayer struct {
	*player
}

func NewHumanPlayer(color PegColor) Player {
	qs := []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "What is your name?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}

	answers := struct{ Name string }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil
	}

	name := answers.Name

	return &humanPlayer{
		player: newPlayer(name, color),
	}
}

func (p *humanPlayer) Shuffle() {
	for cont := true; cont; {
		prompt := &survey.Confirm{
			Message: "You're the dealer. Continue Shuffling?",
			Default: true,
		}

		survey.AskOne(prompt, &cont)
	}
}

func (p *humanPlayer) AddToCrib() []cards.Card {
	cardChoices := make([]string, 0, len(p.hand))
	for _, c := range p.hand {
		cardChoices = append(cardChoices, c.String())
	}

	correctCountValidator := func(val interface{}) error {
		if slice, ok := val.([]string); ok {
			if len(slice) > 2 {
				return fmt.Errorf(`cannot accept a slice with more than length 2 (had length %d)`, len(slice))
			}
		} else if slice, ok := val.([]survey.OptionAnswer); ok {
			if len(slice) > 2 {
				return fmt.Errorf(`cannot accept a slice with more than length 2 (had length %d)`, len(slice))
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("bad type! %T", val)
		}

		// the input is fine
		return nil
	}

	cribCards := []string{}
	prompt := &survey.MultiSelect{
		Message: "Which cards to place in the crib?",
		Options: cardChoices,
	}
	survey.AskOne(prompt, &cribCards, survey.WithValidator(correctCountValidator))

	if len(cribCards) > 2 {
		println(`bad time! never choose more than 2 cards`)
	}
	for i, c := range p.hand {
		for _, cc := range cribCards {
			if c.String() == cc {
				// TODO verify this actually works
				p.hand = append(p.hand[0:i], p.hand[i+1:]...)
			}
		}
	}

	crib := make([]cards.Card, len(cribCards))
	for i, cc := range cribCards {
		crib[i] = cards.NewCardFromString(cc)
	}
	return crib
}

func (p *humanPlayer) Cut() float64 {
	// TODO Ask the user how far down the deck they wanna cut
	return 0
}

func (p *humanPlayer) Peg(maxVal int) (cards.Card, bool, bool) {
	// TODO ask the user which of their cards they would like to peg
	// TODO validate they can peg that
	return cards.Card{}, false, true
}
