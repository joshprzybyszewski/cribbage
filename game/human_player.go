package game

import (
	"fmt"
	"math/rand"
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var _ PlayerInteraction = (*terminalInteraction)(nil)

type terminalInteraction struct {
	myColor PegColor

	scoresByColor   map[PegColor]int
	lagScoreByColor map[PegColor]int
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
	ti := terminalInteraction{
		scoresByColor:   map[PegColor]int{},
		lagScoreByColor: map[PegColor]int{},
	}

	return newPlayer(&ti, name, color)
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
	p.printCurrentScore()
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
		Message: msg + "Which cards to place in the crib?",
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

func (p *terminalInteraction) TellAboutCut(c cards.Card) {
	fmt.Printf("Card cut: %s\n", c.String())
}

func (p *terminalInteraction) AskToPeg(hand, prevPegs []cards.Card, curPeg int) (cards.Card, bool) {
	p.printCurrentScore()
	pegChoices := make([]string, 0, len(hand)+1)
	const sayGo = `Say Go!`
	pegChoices = append(pegChoices, sayGo)
	for _, c := range hand {
		pegChoices = append(pegChoices, c.String())
	}

	canPeg := func(val interface{}) error {
		if oa, ok := val.(survey.OptionAnswer); ok {
			maxValToPeg := maxPeggingValue - curPeg
			if oa.Value == sayGo {
				for _, c := range hand {
					if c.PegValue() <= maxValToPeg {
						return fmt.Errorf("You cannot say go when you have cards to peg")
					}
				}
			} else {
				c := cards.NewCardFromString(oa.Value)
				if c.PegValue() > maxValToPeg {
					return fmt.Errorf("exceeds max peg value: %v", c.String())
				}

			}
			return nil
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("bad type! %T", val)
		}
	}

	msg := `Pegging at: ` + strconv.Itoa(curPeg) + `. The last cards peggged were: `
	for i, c := range prevPegs {
		msg += c.String()
		if i < len(prevPegs)-1 {
			msg += `, `
		} else {
			msg += `. `
		}
	}

	pegCard := ``
	prompt := &survey.Select{
		Message: msg + "Which card to peg next?",
		Options: pegChoices,
	}
	survey.AskOne(prompt, &pegCard, survey.WithValidator(survey.Required), survey.WithValidator(canPeg))

	if pegCard == sayGo {
		return cards.Card{}, true
	}

	return cards.NewCardFromString(pegCard), false
}

func (p *terminalInteraction) TellAboutScores(cur, lag map[PegColor]int, msgs ...string) {
	for c, s := range cur {
		if n := s - p.scoresByColor[c]; n != 0 {
			if c == p.myColor {
				fmt.Printf("You scored %d points for %v\n", n, msgs)
			} else {
				fmt.Printf("%s scored %d points for %v\n", c.String(), n, msgs)
			}
		}
		p.scoresByColor[c] = cur[c]
		p.lagScoreByColor[c] = lag[c]
	}
}

func (p *terminalInteraction) printCurrentScore() {
	if len(p.scoresByColor) != 2 {
		return
	}
	fmt.Println(`------------`)
	fmt.Printf("  You: %3d\n", p.scoresByColor[p.myColor])
	for c, s := range p.scoresByColor {
		if c != p.myColor {
			fmt.Printf("%5s: %3d\n", c.String(), s)
		}
	}
	fmt.Println(`------------`)
}
