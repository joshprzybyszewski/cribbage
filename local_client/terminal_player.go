package local_client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	serverDomain = `http://localhost:8080`
)

type terminalClient struct {
	server *http.Client

	me            model.Player
	myCurrentGame model.GameID
	myGames       map[model.GameID]model.Game
}

type terminalRequest struct {
	gameID model.GameID
	game model.Game
	msg string
}

func StartTerminalInteraction() error {
	tc := terminalClient{
		server:  &http.Client{},
		myGames: make(map[model.GameID]model.Game),
	}
	if tc.shouldSignIn() {
		tc.me.ID = tc.getPlayerID(`What is your player ID?`)
	} else {
		err := tc.createPlayer()
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	
	port := 8081 + (int(tc.me.ID) % 100)
	reqChan := make(chan terminalRequest, 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		filename := fmt.Sprintf("./playerlogs/%d.log", tc.me.ID)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("failed opening file: %s", err)
			return
		}
		defer f.Close()

		playerServerFile := bufio.NewWriter(f)

		router := gin.New()
		router.Use(gin.LoggerWithWriter(playerServerFile), gin.Recovery())

		router.POST("/blocking/:gameID", func(c *gin.Context) {
			gIDStr := c.Param("gameID")
			n, err := strconv.Atoi(gIDStr)
			if err != nil {
				c.String(http.StatusBadRequest, "Invalid GameID: %s", gIDStr)
				return
			}

			reqChan <- terminalRequest{
				gameID: model.GameID(n),
				msg: `you're blocking`,
			}

			c.String(http.StatusOK, `received`)
		})
		router.POST("/message/:gameID", func(c *gin.Context) {
			gIDStr := c.Param("gameID")
			n, err := strconv.Atoi(gIDStr)
			if err != nil {
				c.String(http.StatusBadRequest, "Invalid GameID: %s", gIDStr)
				return
			}

			reqChan <- terminalRequest{
				gameID: model.GameID(n),
				msg: `Received Message`,
			}
			c.String(http.StatusOK, `received`)
		})
		router.POST("/score/:gameID", func(c *gin.Context) {
			gIDStr := c.Param("gameID")
			n, err := strconv.Atoi(gIDStr)
			if err != nil {
				c.String(http.StatusBadRequest, "Invalid GameID: %s", gIDStr)
				return
			}

			reqChan <- terminalRequest{
				gameID: model.GameID(n),
				msg: `Received score update`,
			}
			c.String(http.StatusOK, `received`)
		})

		router.Run(fmt.Sprintf("0.0.0.0:%d", port)) // listen and serve on the addr
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Let the server know about where we're serving our listener
		url := fmt.Sprintf("/create/interaction/%d/localhost/%d", tc.me.ID, port)
		tc.makeRequest(`POST`, url, nil)
	}()

	wg.Add(1)
	go func() {
		if tc.shouldCreateGame() {
			err := tc.createGame()
			if err != nil {
				return
			}
		}

		err := tc.updatePlayer()
		if err != nil {
			return
		}

		reqChan <- terminalRequest{
			game: tc.myGames[tc.myCurrentGame],
		}
		for {
			select {
			case req := <-reqChan:
				fmt.Printf("Message: \"%s\"\n", req.msg)
				gID := req.gameID
				if req.gameID == model.InvalidGameID {
					gID = req.game.ID
				}
				err := tc.requestAndSendAction(gID)
				if err != nil {
					reqChan <- terminalRequest{
						gameID: gID,
						msg: `Problem doing action. Try again?`,
					}
				}
			}
		}
	}()

	// Block until forever...?
	wg.Wait()

	return nil
}

func (tc *terminalClient) makeRequest(method, apiURL string, data io.Reader) ([]byte, error) {
	urlStr := serverDomain + apiURL
	req, err := http.NewRequest(method, urlStr, data)
	if err != nil {
		return nil, err
	}

	response, err := tc.server.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %+v\n%s", response, response.Body)
	}

	return ioutil.ReadAll(response.Body)
}

func (tc *terminalClient) createPlayer() error {
	name := tc.getName()

	respBytes, err := tc.makeRequest(`POST`, `/create/player/`+name, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, &tc.me)
	if err != nil {
		return err
	}

	fmt.Printf("Your player ID is: %v\n", tc.me.ID)

	return nil
}

func (p *terminalClient) shouldSignIn() bool {
	should := true

	prompt := &survey.Confirm{
		Message: "Sign in?",
		Default: true,
	}

	survey.AskOne(prompt, &should)
	return should
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

func (p *terminalClient) shouldCreateGame() bool {
	cont := true

	prompt := &survey.Confirm{
		Message: "Create new game?",
		Default: true,
	}

	survey.AskOne(prompt, &cont)
	return cont
}

func (tc *terminalClient) createGame() error {
	opID := tc.getPlayerID("What's your opponent's ID?")
	url := fmt.Sprintf("/create/game/%d/%v", opID, tc.me.ID)

	respBytes, err := tc.makeRequest(`POST`, url, nil)
	if err != nil {
		return err
	}

	g := model.Game{}

	err = json.Unmarshal(respBytes, &g)
	if err != nil {
		return err
	}

	tc.myCurrentGame = g.ID
	tc.myGames[g.ID] = g

	return nil
}

func (tc *terminalClient) getPlayerID(msg string) model.PlayerID {
	qs := []*survey.Question{{
		Name:      "id",
		Prompt:    &survey.Input{Message: msg},
		Validate:  survey.Required,
		Transform: survey.Title,
	}}

	answers := struct{ Id int }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.InvalidPlayerID
	}

	return model.PlayerID(answers.Id)
}

func (tc *terminalClient) updatePlayer() error {
	url := fmt.Sprintf("/player/%v", tc.me.ID)

	respBytes, err := tc.makeRequest(`GET`, url, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, &tc.me)
	if err != nil {
		return err
	}

	for gID := range tc.me.Games {
		g, err := tc.requestGame(gID)
		if err != nil {
			return err
		}

		tc.myGames[gID] = g

		if !g.IsOver() {
			tc.myCurrentGame = gID
		}
	}

	return nil
}

func (tc *terminalClient) getGame(gID model.GameID) (model.Game, error) {
	g, ok := tc.myGames[gID]
	if ok {
		return g, nil
	}

	return tc.requestGame(gID)
}

func (tc *terminalClient) requestGame(gID model.GameID) (model.Game, error) {
	url := fmt.Sprintf("/game/%v", gID)

	respBytes, err := tc.makeRequest(`GET`, url, nil)
	if err != nil {
		return model.Game{}, err
	}

	g := model.Game{}

	err = json.Unmarshal(respBytes, &g)
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
}
