package local_client

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"

	"github.com/joshprzybyszewski/cribbage/model"
)

type terminalClient struct {
	server *http.Client

	myID          model.PlayerID
	myCurrentGame model.GameID
	myGames       map[model.GameID]model.Game
}

func StartTerminalInteraction() error {
	tc := terminalClient{
		server: &http.Client{},
	}
	err := tc.createPlayer()
	if err != nil {
		return err
	}

	return nil
}

func (tc *terminalClient) makeRequest(method, apiURL string, data io.ReadCloser) error {
	url, err := url.Parse(apiURL)
	if err != nil {
		return err
	}

	request := http.Request{
		Method: method,
		URL:    url,
		Body:   data,
	}
	response, err := tc.server.Do(&request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (tc *terminalClient) createPlayer() error {
	name := tc.getName()

	err := tc.makeRequest(`POST`, `/create/player/`+name, nil)
	if err != nil {
		return err
	}

	// TODO get player id and start storing that

	return nil
}

func (tc *terminalClient) getName() string {
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
		return `Player`
	}

	return answers.Name
}

func (tc *terminalClient) getPlayerAction() (model.PlayerAction, error) {
	g, ok := tc.myGames[tc.myCurrentGame]
	if !ok {
		return model.PlayerAction{}, errors.New(`does not have game to play`)
	}

	b, ok := g.BlockingPlayers[tc.myID]
	if !ok {
		return model.PlayerAction{}, errors.New(`I am not blocking play right now`)
	}

	tc.printCurrentScore()

	var action interface{}
	switch b {
	case model.DealCards:
		action = tc.getDealAction()
	case model.CribCard:
		action = tc.getBuildCribAction()
	case model.CutCard:
		action = tc.getCutDeckAction()
	case model.PegCard:
		action = tc.getPegAction()
	case model.CountHand:
		action = tc.getCountHandAction()
	case model.CountCrib:
		action = tc.getCountCribAction()
	}

	return model.PlayerAction{
		GameID:    tc.myCurrentGame,
		ID:        tc.myID,
		Overcomes: b,
		Action:    action,
	}, nil
}

func (p *terminalClient) getDealAction() model.DealAction {
	i := 1
	cont := true
	for cont {
		msg := fmt.Sprintf("You've shuffled %d times. Continue?", i)
		prompt := &survey.Confirm{
			Message: msg,
			Default: true,
		}

		survey.AskOne(prompt, &cont)
	}

	return model.DealAction{
		NumShuffles: i,
	}
}

func (tc *terminalClient) getBuildCribAction() model.BuildCribAction {
	g := tc.myGames[tc.myCurrentGame]
	myHand := g.Hands[tc.myID]
	desired := len(myHand) - 4
	cardChoices := make([]string, 0, len(myHand))
	for _, c := range myHand {
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

	msg := ``
	if g.CurrentDealer == tc.myID {
		msg = `Your crib.`
	} else {
		myColor, dealerColor := g.PlayerColors[tc.myID], g.PlayerColors[g.CurrentDealer]
		if myColor == dealerColor {
			msg = `Partner's crib.`
		} else {
			msg = fmt.Sprintf(`%s's crib.`, dealerColor.String())
		}
	}

	cribCards := []string{}
	prompt := &survey.MultiSelect{
		Message: msg + " Which cards to place in the crib?",
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
func (tc *terminalClient) getPegAction() model.PegAction {
	g := tc.myGames[tc.myCurrentGame]
	myHand := g.Hands[tc.myID]

	pegChoices := make([]string, 0, len(myHand)+1)
	const sayGo = `Say Go!`
	pegChoices = append(pegChoices, sayGo)
	for _, c := range myHand {
		hasPegged := false
		for _, pc := range g.PeggedCards {
			if pc.Card == c {
				hasPegged = true
				break
			}
		}
		if !hasPegged {
			pegChoices = append(pegChoices, c.String())
		}
	}

	curPeg := g.CurrentPeg()
	maxValToPeg := model.MaxPeggingValue - curPeg

	canPeg := func(val interface{}) error {
		if oa, ok := val.(survey.OptionAnswer); ok {
			if oa.Value == sayGo {
				for _, c := range myHand {
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
			msg += `.`
		}
	}

	pegCardStr := ``
	prompt := &survey.Select{
		Message: msg + " Which card to peg next?",
		Options: pegChoices,
		Filter:  func(filter string, value string, index int) bool { return true },
	}
	survey.AskOne(prompt, &pegCardStr, survey.WithValidator(survey.Required), survey.WithValidator(canPeg))

	if pegCardStr == sayGo {
		return model.PegAction{
			Card:  model.Card{},
			SayGo: true,
		}
	}

	return model.PegAction{
		Card:  model.NewCardFromString(pegCardStr),
		SayGo: false,
	}
}
func (tc *terminalClient) getCountHandAction() model.CountHandAction {
	g := tc.myGames[tc.myCurrentGame]
	myHand := g.Hands[tc.myID]

	msg := fmt.Sprintf("Cut: %s, Hand: (%s %s %s %s).",
		g.CutCard,
		myHand[0],
		myHand[1],
		myHand[2],
		myHand[3],
	)

	qs := []*survey.Question{
		{
			Name:      "pts",
			Prompt:    &survey.Input{Message: msg + " How many points in your hand?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}

	answers := struct {
		HandPts int `survey:"pts"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.CountHandAction{
			Pts: 0,
		}
	}

	return model.CountHandAction{
		Pts: answers.HandPts,
	}
}
func (tc *terminalClient) getCountCribAction() model.CountCribAction {
	// TODO validate that I am the current dealer?
	g := tc.myGames[tc.myCurrentGame]

	msg := fmt.Sprintf("Cut: %s, Hand: (%s %s %s %s).",
		g.CutCard,
		g.Crib[0],
		g.Crib[1],
		g.Crib[2],
		g.Crib[3],
	)

	qs := []*survey.Question{
		{
			Name:      "pts",
			Prompt:    &survey.Input{Message: msg + " How many points in the crib?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}

	answers := struct {
		CribPts int `survey:"pts"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.CountCribAction{
			Pts: 0,
		}
	}

	return model.CountCribAction{
		Pts: answers.CribPts,
	}
}

func (tc *terminalClient) printCurrentScore() {
	g := tc.myGames[tc.myCurrentGame]
	myColor := g.PlayerColors[tc.myID]
	fmt.Println(`----------`)
	fmt.Printf("%5s (you): %3d\n", myColor.String(), g.CurrentScores[myColor])
	for c, s := range g.CurrentScores {
		if c != myColor {
			fmt.Printf("%11s: %3d\n", c.String(), s)
		}
	}
	fmt.Println(`----------`)
}

// TODO how will we tell about the cut?
// func (p *terminalClient) TellAboutCut(c model.Card) {
// 	fmt.Printf("Card cut: %s\n", c.String())
// }

// TODO how will we notify of score updates too?
// func (p *terminalClient) TellAboutScores(cur, lag map[model.PlayerColor]int, msgs ...string) {
// 	for c, s := range cur {
// 		if n := s - p.scoresByColor[c]; n != 0 {
// 			if c == p.myColor {
// 				fmt.Printf("You   scored %2d points for %v\n", n, msgs)
// 			} else {
// 				fmt.Printf("%-5s scored %2d points for %v\n", c.String(), n, msgs)
// 			}
// 		}
// 		p.scoresByColor[c] = cur[c]
// 		p.lagScoreByColor[c] = lag[c]
// 	}
// }
