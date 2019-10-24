package main

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/strategy"
)

func main() {
	for {
		qs := []*survey.Question{
			{
				Name:      "inputHand",
				Prompt:    &survey.Input{Message: "What is your dealt hand?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
		}
		var inputHand string
		err := survey.Ask(qs, &inputHand)
		if err != nil {
			fmt.Printf("got error: %+v\n", err)
		}
		if inputHand == `exit` || inputHand == `quit` {
			os.Exit(0)
		}
		inputCards := strings.Split(inputHand, `,`)

		reportAboutHand(inputCards)
	}
}

func reportAboutHand(cstrs []string) {
	fmt.Printf("Calculating for hand: %+v\n", strToCards(cstrs))
	lowCrib := strategy.GiveCribLowestPotential(0, strToCards(cstrs))
	fmt.Printf("GiveCribLowestPotential: %+v\n", lowCrib)

	highCrib := strategy.GiveCribHighestPotential(0, strToCards(cstrs))
	fmt.Printf("GiveCribHighestPotential: %+v\n", highCrib)
}

func strToCards(s []string) []cards.Card {
	c := make([]cards.Card, len(s))
	for i, str := range s {
		c[i] = cards.NewCardFromString(str)
	}
	return c
}
