package game

import (
	"fmt"
	"math/rand"

	"github.com/AlecAivazis/survey"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var _ PlayerInteraction = (*terminalInteraction)(nil)

type terminalInteraction struct {
	myColor PegColor
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

	return newPlayer(&terminalInteraction{}, name, color)
}

func (p *terminalInteraction) AskToShuffle() bool {
	cont := true
		prompt := &survey.Confirm{
			Message: "You're the dealer. Continue Shuffling?",
			Default: true,
		}

		survey.AskOne(prompt, &cont)
		return cont
}

func (p *terminalInteraction) AskForCribCards(dealerColor PegColor, desired int, hand []cards.Card) []cards.Card {
	cardChoices := make([]string, 0, len(hand))
	for _, c := range hand {
		cardChoices = append(cardChoices, c.String())
	}

	correctCountValidator := func(val interface{}) error {
		if slice, ok := val.([]string); ok {
			if len(slice) != desired {
				return fmt.Errorf(`cannot accept a slice with more than length %d (had length %d)`, desired, len(slice))
			}
		} else if slice, ok := val.([]survey.OptionAnswer); ok {
			if len(slice) != desired {
				return fmt.Errorf(`cannot accept a slice with more than length %d (had length %d)`, desired, len(slice))
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("bad type! %T", val)
		}

		// the input is fine
		return nil
	}

	msg := `Crib goes to you. `
	if p.myColor != dealerColor {
		msg = `Crib does not go to you. `
	}
	cribCards := []string{}
	prompt := &survey.MultiSelect{
		Message: msg+"Which cards to place in the crib?",
		Options: cardChoices,
	}
	survey.AskOne(prompt, &cribCards, survey.WithValidator(correctCountValidator))

	if len(cribCards) != desired {
		println(`bad time! never choose more than 2 cards`)
		return nil
	}

	crib := make([]cards.Card, len(cribCards))
	for i, cc := range cribCards {
		crib[i] = cards.NewCardFromString(cc)
	}
	return crib
}

func (p *terminalInteraction) AskForCut() float64 {
	const thin = `thin`
	const middle = `middle`
	const thick = `thick`
	cutChoice := ``
	prompt := &survey.Select{
		Message: "How would you like to cut?",
		Options: []string{thin, middle, thick},
	}
	survey.AskOne(prompt, &cutChoice)

	switch cutChoice {
	case thin:
		return (rand.Float64() + 0) / 3
	case middle:
		return (rand.Float64() + 1) / 3
	case thick:
		return (rand.Float64() + 2) / 3
	}
	
	return 0.500
}

func (p *terminalInteraction) AskToPeg(hand, prevPegs []cards.Card, curPeg int) (cards.Card, bool) {
	// TODO ask the user which of their cards they would like to peg
	cardChoices := make([]string, 0, len(hand))
	for _, c := range hand {
		cardChoices = append(cardChoices, c.String())
	}

	msg := `The last `

	cribCards := []string{}
	prompt := &survey.Select{
		Message: msg+" Which card to peg?",
		Options: cardChoices,
	}
	survey.AskOne(prompt, &cribCards, survey.WithValidator(correctCountValidator))

	if len(cribCards) != desired {
		println(`bad time! never choose more than 2 cards`)
		return nil
	}

	crib := make([]cards.Card, len(cribCards))
	for i, cc := range cribCards {
		crib[i] = cards.NewCardFromString(cc)
	}
	return crib

	return cards.Card{}, false
}
