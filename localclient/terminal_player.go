package localclient

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

type termReqType int

const (
	blocking termReqType = iota
	scoreUpdate
	message
)

type terminalRequest struct {
	gameID model.GameID
	game   model.Game
	req    termReqType
	msg    string
}

func StartTerminalInteraction() error {
	tc := terminalClient{
		server:  &http.Client{},
		myGames: make(map[model.GameID]model.Game),
	}
	if tc.shouldSignIn() {
		tc.me.ID = tc.getPlayerID(`What is your username?`)
	} else {
		err := tc.createPlayer()
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	port := 8081 + (len(tc.me.ID) % 100)
	reqChan := make(chan terminalRequest, 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		dir := `./playerlogs`
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("failed creating directory \"%s\": %v\n", dir, err)
				return
			}
		}
		filename := fmt.Sprintf(dir+"/%d.log", tc.me.ID)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("failed opening file: %v", err)
			return
		}
		defer f.Close()

		playerServerFile := bufio.NewWriter(f)

		router := gin.New()
		router.Use(gin.LoggerWithWriter(playerServerFile), gin.Recovery())

		router.POST("/blocking/:gameID", handleBlocking(reqChan))
		router.POST("/message/:gameID", handleMessage(reqChan))
		router.POST("/score/:gameID", handleScoreUpdate(reqChan))

		err = router.Run(fmt.Sprintf("0.0.0.0:%d", port)) // listen and serve on the addr
		fmt.Printf("router.Run error: %+v\n", err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Let the server know about where we're serving our listener
		url := fmt.Sprintf("/create/interaction/%s/localhost/%d", tc.me.ID, port)
		_, err := tc.makeRequest(`POST`, url, nil)
		if err != nil {
			fmt.Printf("Error telling server about interaction: %+v\n", err)
		}
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
			msg:  `Starting terminal player`,
			req:  message,
		}
		for req := range reqChan {
			fmt.Printf("Message: \"%s\"\n", req.msg)
			gID := req.gameID
			if req.gameID == model.InvalidGameID {
				gID = req.game.ID
			}
			if gID == model.InvalidGameID {
				continue
			}
			switch req.req {
			case blocking:
				err := tc.requestAndSendAction(gID)
				if err != nil {
					reqChan <- terminalRequest{
						gameID: gID,
						msg:    `Problem doing action. Try again?`,
						req:    message,
					}
				}
			case message:
				fmt.Println(req.msg)
			case scoreUpdate:
				fmt.Println(req.msg)
				tc.printCurrentScore()
			default:
				fmt.Printf("Developer error: req needs a type %+v\n", req)
			}
		}
	}()

	// Block until forever...?
	wg.Wait()

	return nil
}
func handleBlocking(reqChan chan terminalRequest) func(*gin.Context) {
	return func(c *gin.Context) {
		gID, msg, err := getGameIDAndBody(c, `We heard you're blocking`)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		reqChan <- terminalRequest{
			gameID: gID,
			msg:    msg,
			req:    blocking,
		}

		c.String(http.StatusOK, `received`)
	}
}

func handleMessage(reqChan chan terminalRequest) func(*gin.Context) {
	return func(c *gin.Context) {
		gID, msg, err := getGameIDAndBody(c, `Received a message`)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		reqChan <- terminalRequest{
			gameID: gID,
			msg:    msg,
			req:    message,
		}
		c.String(http.StatusOK, `received`)
	}
}

func handleScoreUpdate(reqChan chan terminalRequest) func(*gin.Context) {
	return func(c *gin.Context) {
		gID, msg, err := getGameIDAndBody(c, `There was a score update`)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		reqChan <- terminalRequest{
			gameID: gID,
			msg:    msg,
			req:    scoreUpdate,
		}
		c.String(http.StatusOK, `received`)
	}
}

func getGameIDAndBody(c *gin.Context, defBody string) (model.GameID, string, error) {
	gIDStr := c.Param("gameID")
	n, err := strconv.Atoi(gIDStr)
	if err != nil {
		return model.InvalidGameID, ``, fmt.Errorf("Invalid GameID: %s", gIDStr)
	}

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	var msg string
	if err != nil || len(reqBody) == 0 {
		msg = defBody
	} else {
		msg = string(reqBody)
	}

	return model.GameID(n), msg, nil
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
	username, name := tc.getName()

	respBytes, err := tc.makeRequest(`POST`, `/create/player/`+username+`/`+name, nil)
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

func (tc *terminalClient) shouldSignIn() bool {
	should := true

	prompt := &survey.Confirm{
		Message: "Sign in?",
		Default: true,
	}

	err := survey.AskOne(prompt, &should)
	if err != nil {
		fmt.Printf("survey.AskOne error: %+v\n", err)
		return false
	}
	return should
}

func (tc *terminalClient) getName() (string, string) {
	qs := []*survey.Question{
		{
			Name:      "username",
			Prompt:    &survey.Input{Message: "What username do you want?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "What is your name?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}

	answers := struct{ Username, Name string }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return `username`, `Player`
	}

	return answers.Username, answers.Name
}

func (tc *terminalClient) shouldCreateGame() bool {
	cont := true

	prompt := &survey.Confirm{
		Message: "Create new game?",
		Default: true,
	}

	err := survey.AskOne(prompt, &cont)
	if err != nil {
		fmt.Printf("survey.AskOne error: %+v\n", err)
		return false
	}
	return cont
}

func (tc *terminalClient) createGame() error {
	opID := tc.getPlayerID("What's your opponent's username?")
	url := fmt.Sprintf("/create/game/%s/%s", opID, tc.me.ID)

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
		Name:      "username",
		Prompt:    &survey.Input{Message: msg},
		Validate:  survey.Required,
		Transform: survey.Title,
	}}

	answers := struct{ Username string }{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return model.InvalidPlayerID
	}

	return model.PlayerID(answers.Username)
}

func (tc *terminalClient) updatePlayer() error {
	url := fmt.Sprintf("/player/%s", tc.me.ID)

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
	tc.myGames[g.ID] = g

	return g, nil
}
