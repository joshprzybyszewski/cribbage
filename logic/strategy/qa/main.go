package main

import (
	"fmt"
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
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
		if inputHand == `Exit` || inputHand == `Quit` || inputHand == `Q` {
			os.Exit(0)
		}
		inputCards := strings.Split(inputHand, `,`)

		reportAboutHand(inputCards)
	}
}

func reportAboutHand(cstrs []string) {
	fmt.Printf("Calculating for hand: %+v\n", strToCards(cstrs))
	lowCrib, err := strategy.GiveCribLowestPotential(0, strToCards(cstrs))
	if err != nil {
		fmt.Printf("GiveCribLowestPotential: Error! %v\n", err)
	} else {
		fmt.Printf("GiveCribLowestPotential: %+v\n", lowCrib)
	}

	highCrib, err := strategy.GiveCribHighestPotential(0, strToCards(cstrs))
	if err != nil {
		fmt.Printf("GiveCribHighestPotential: Error! %v\n", err)
	} else {
		fmt.Printf("GiveCribHighestPotential: %+v\n", highCrib)
	}
}

func strToCards(s []string) []model.Card {
	c := make([]model.Card, len(s))
	for i, str := range s {
		c[i] = model.NewCardFromString(str)
	}
	return c
}
