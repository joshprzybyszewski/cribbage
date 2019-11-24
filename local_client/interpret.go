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
		// c.String(http.StatusBadRequest, "Invalid GameID: %s", gIDStr)
		return err
	}
	b, err := json.Marshal(pa)
	if err != nil {
		// c.String(http.StatusBadRequest, "Bad Marshaling: %s", gIDStr)
		return err
	}
	buf := bytes.NewBuffer(b)

	url := fmt.Sprintf("/action/%d", g.ID)
	bytes, err := tc.makeRequest(`POST`, url, buf)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		// c.String(http.StatusBadRequest, "Bad Marshaling: %s", gIDStr)
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
