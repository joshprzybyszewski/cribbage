package local_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/joshprzybyszewski/cribbage/model"
)

var (
	errIAmNotBlocking = errors.New(`i'm not blocking`)
)

func (tc *terminalClient) requestAndSendAction(gID model.GameID) error {
	// TODO change this to a "get game" call because we don't want to make network
	// requests all the time
	g, err := tc.requestGame(gID)
	if err != nil {
		return err
	}

	pa, err := tc.getPlayerAction(g)
	if err != nil {
		if err == errIAmNotBlocking {
			fmt.Printf("Waiting for other players\n")
			return nil
		}
		return err
	}
	b, err := json.Marshal(pa)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)

	url := fmt.Sprintf("/action/%d", g.ID)
	_, err = tc.makeRequest(`POST`, url, buf)
	return err
}


func (tc *terminalClient) getPlayerAction(g model.Game) (model.PlayerAction, error) {
	b, ok := g.BlockingPlayers[tc.me.ID]
	if !ok {
		return model.PlayerAction{}, errIAmNotBlocking
	}

	tc.printCurrentScore()

	var action interface{}
	switch b {
	case model.DealCards:
		action = tc.getDealAction()
	case model.CribCard:
		action = tc.getBuildCribAction(g)
	case model.CutCard:
		action = tc.getCutDeckAction()
	case model.PegCard:
		action = tc.getPegAction(g)
	case model.CountHand:
		action = tc.getCountHandAction(g)
	case model.CountCrib:
		action = tc.getCountCribAction(g)
	}

	return model.PlayerAction{
		GameID:    g.ID,
		ID:        tc.me.ID,
		Overcomes: b,
		Action:    action,
	}, nil
}

func (p *terminalClient) getDealAction() model.DealAction {
	qs := []*survey.Question{{
		Name:      "numShuffles",
		Prompt:    &survey.Input{Message: `How many times to shuffle?`},
		Validate:  survey.Required,
		Transform: survey.Title,
	}}

	answers := struct{ NumShuffles int }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.DealAction{}
	}

	return model.DealAction{
		NumShuffles: answers.NumShuffles,
	}
}

func (tc *terminalClient) getBuildCribAction(g model.Game) model.BuildCribAction {
	hand := g.Hands[tc.me.ID]
	desired := len(hand) - 4
	cardChoices := make([]string, 0, len(hand))
	for _, c := range hand {
		cardChoices = append(cardChoices, c.String())
	}

	correctCountValidator := func(val interface{}) error {
		if slice, ok := val.([]string); ok {
			if len(slice) != desired {
				return fmt.Errorf(`Asked for %d cards. (You gave us %d)`, desired, len(slice))
			}
		} else if slice, ok := val.([]survey.OptionAnswer); ok {
			if len(slice) != desired {
				return fmt.Errorf(`Asked for %d cards. (You gave us %d)`, desired, len(slice))
			}
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("bad type! %T", val)
		}

		// the input is fine
		return nil
	}

	msg := `Crib does not go to you. `
	if tc.me.ID == g.CurrentDealer {
		msg = `Crib goes to you. `
	} else if g.PlayerColors[tc.me.ID] == g.PlayerColors[g.CurrentDealer] {
		msg = `Crib goes to your partner. `
	}

	cribCards := []string{}
	prompt := &survey.MultiSelect{
		Message: msg + "Which cards to place in the crib?",
		Options: cardChoices,
	}
	survey.AskOne(prompt, &cribCards, survey.WithValidator(correctCountValidator))

	if len(cribCards) != desired {
		fmt.Printf(`bad time! expected %d cards, received %d`, desired, len(cribCards))
		return model.BuildCribAction{}
	}

	crib := make([]model.Card, len(cribCards))
	for i, cc := range cribCards {
		crib[i] = model.NewCardFromString(cc)
	}

	return model.BuildCribAction{
		Cards: crib,
	}
}
func (tc *terminalClient) getCutDeckAction() model.CutDeckAction {
	const thin = `thin`
	const middle = `middle`
	const thick = `thick`
	cutChoice := ``
	prompt := &survey.Select{
		Message: "How would you like to cut?",
		Options: []string{thin, middle, thick},
		Filter:  func(filter string, value string, index int) bool { return true },
	}
	survey.AskOne(prompt, &cutChoice)

	perc := 0.500
	switch cutChoice {
	case thin:
		perc = (rand.Float64() + 0) / 3
	case middle:
		perc = (rand.Float64() + 1) / 3
	case thick:
		perc = (rand.Float64() + 2) / 3
	}

	return model.CutDeckAction{
		Percentage: perc,
	}
}
func (tc *terminalClient) getPegAction(g model.Game) model.PegAction {
	hand := g.Hands[tc.me.ID]
	curPeg := g.CurrentPeg()
	
	pegChoices := make([]string, 0, len(hand)+1)
	const sayGoOption = `Say Go!`
	pegChoices = append(pegChoices, sayGoOption)
	for _, c := range hand {
		pegChoices = append(pegChoices, c.String())
	}

	canPeg := func(val interface{}) error {
		if oa, ok := val.(survey.OptionAnswer); ok {
			maxValToPeg := model.MaxPeggingValue - curPeg
			if oa.Value == sayGoOption {
				for _, c := range hand {
					if c.PegValue() <= maxValToPeg {
						return fmt.Errorf("You cannot say go when you have cards to peg")
					}
				}
			} else {
				c := model.NewCardFromString(oa.Value)
				if c.PegValue() > maxValToPeg {
					return fmt.Errorf("Card (%v) exceeds max peg value (%d)", c.String(), maxValToPeg)
				}

			}
			return nil
		} else {
			// otherwise we cannot convert the value into a string and cannot enforce length
			return fmt.Errorf("bad type! %T", val)
		}
	}

	msg := `Pegging at: ` + strconv.Itoa(curPeg) + `. The last cards pegged were: `
	for i, c := range g.PeggedCards {
		msg += c.String()
		if i < len(g.PeggedCards)-1 {
			msg += `, `
		} else {
			msg += `. `
		}
	}

	pegCard := ``
	prompt := &survey.Select{
		Message: msg + "Which card to peg next?",
		Options: pegChoices,
		Filter:  func(filter string, value string, index int) bool { return true },
	}
	survey.AskOne(prompt, &pegCard, survey.WithValidator(survey.Required), survey.WithValidator(canPeg))

	c := model.Card{}
	sayGo := pegCard == sayGoOption
	if !sayGo {
		c = model.NewCardFromString(pegCard)
	}

	return model.PegAction{
			Card: c,
			SayGo: sayGo,
		}
}
func (tc *terminalClient) getCountHandAction(g model.Game) model.CountHandAction {
	hand := g.Hands[tc.me.ID]

	msg := fmt.Sprintf(`Cut Card: %s, Hand: `, g.CutCard)
	for i, c := range hand {
		msg += c.String()
		if i < len(g.PeggedCards)-1 {
			msg += `, `
		} else {
			msg += `. `
		}
	}

	qs := []*survey.Question{{
		Name:      "handPoints",
		Prompt:    &survey.Input{Message: msg + `How many points in your hand?`},
		Validate:  survey.Required,
		Transform: survey.Title,
	}}

	answers := struct{ HandPoints int }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.CountHandAction{}
	}

	return model.CountHandAction{
		Pts: answers.HandPoints,
	}
}
func (tc *terminalClient) getCountCribAction(g model.Game) model.CountCribAction {
	crib := g.Crib

	msg := fmt.Sprintf(`Cut Card: %s, Crib: `, g.CutCard)
	for i, c := range crib {
		msg += c.String()
		if i < len(g.PeggedCards)-1 {
			msg += `, `
		} else {
			msg += `. `
		}
	}

	qs := []*survey.Question{{
		Name:      "cribPoints",
		Prompt:    &survey.Input{Message: msg + `How many points in the crib?`},
		Validate:  survey.Required,
		Transform: survey.Title,
	}}

	answers := struct{ CribPoints int }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.CountCribAction{}
	}

	return model.CountCribAction{
		Pts: answers.CribPoints,
	}
}

func (tc *terminalClient) printCurrentScore() {
	g := tc.myGames[tc.myCurrentGame]
	myColor := g.PlayerColors[tc.me.ID]
	fmt.Println(`----------`)
	fmt.Printf("%5s (you): %3d\n", myColor.String(), g.CurrentScores[myColor])
	for c, s := range g.CurrentScores {
		if c != myColor {
			fmt.Printf("%11s: %3d\n", c.String(), s)
		}
	}
	fmt.Println(`----------`)
}
