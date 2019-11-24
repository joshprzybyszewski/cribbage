package local_client

import (
	"bytes"
	// "bufio"
	"encoding/json"
	"errors"
	"fmt"
	// "io"
	// "io/ioutil"
	// "math/rand"
	// "net/http"
	// "net/url"
	// "os"
	// "strconv"
	// "sync"

	survey "github.com/AlecAivazis/survey/v2"
	// "github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
)

var (
	errIAmNotBlocking = errors.New(`i'm not blocking`)
)

func (tc *terminalClient) requestAndSendAction(gID model.GameID) error {
	g, err := tc.getGame(gID)
	if err != nil {
		return err
	}

	pa, err := tc.askForAction(g)
	if err != nil {
		if err == errIAmNotBlocking {
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
	bytes, err := tc.makeRequest(`POST`, url, buf)
	if err != nil {
		fmt.Printf("err: `%+v` %+v\n", string(bytes), err)
		return err
	}
	fmt.Printf("bytes: %+v\n", string(bytes))


	return nil
}

func (tc *terminalClient) askForAction(g model.Game) (model.PlayerAction, error) {
	r, ok := g.BlockingPlayers[tc.me.ID]
	if !ok {
		fmt.Printf("Waiting...\n")
		return model.PlayerAction{}, errIAmNotBlocking
	}

	switch r {
	case model.DealCards:
		return tc.askForDeal()
	case model.CribCard:
		return tc.askForCrib(g)
	case model.CutCard:
		// return tc.askForDeal()
	case model.PegCard:
		// return tc.askForDeal()
	case model.CountHand:
		// return tc.askForDeal()
	case model.CountCrib:
		// return tc.askForDeal()
	}

	return model.PlayerAction{}, errors.New(`unhandleable state?`)
}

func (tc *terminalClient) askForDeal() (model.PlayerAction, error) {
	qs := []*survey.Question{{
		Name:      "numShuffles",
		Prompt:    &survey.Input{Message: `How many times to shuffle?`},
		Validate:  survey.Required,
		Transform: survey.Title,
	}}

	answers := struct{ NumShuffles int }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.PlayerAction{}, err
	}

	pa := model.PlayerAction{
		GameID:    tc.myCurrentGame,
		ID:        tc.me.ID,
		Overcomes: model.DealCards,
		Action: model.DealAction{
			NumShuffles: answers.NumShuffles,
		},
	}

	return pa, nil
}

func (tc *terminalClient) askForCrib(g model.Game) (model.PlayerAction, error) {
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
		return model.PlayerAction{}, errors.New(`should not hit this`)
	}

	crib := make([]model.Card, len(cribCards))
	for i, cc := range cribCards {
		crib[i] = model.NewCardFromString(cc)
	}

	pa := model.PlayerAction{
		GameID:    tc.myCurrentGame,
		ID:        tc.me.ID,
		Overcomes: model.DealCards,
		Action: model.BuildCribAction{
			Cards: crib,
		},
	}

	return pa, nil
}
