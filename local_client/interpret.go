package local_client

import (
	// "bufio"
	// "encoding/json"
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

func (tc *terminalClient) askForAction(g model.Game) (model.PlayerAction, error) {
	r, ok := g.BlockingPlayers[tc.me.ID]
	if !ok {
		fmt.Printf("Waiting...\n")
		return model.PlayerAction{}, nil
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
		GameID: tc.myCurrentGame,
		ID: tc.me.ID,
		Overcomes: model.DealCards,
		Action: model.DealAction{
			NumShuffles: answers.NumShuffles,
		},
	}

	return pa, nil
}