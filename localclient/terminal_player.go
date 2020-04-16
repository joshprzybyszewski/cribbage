package localclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joshprzybyszewski/cribbage/jsonutils"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	serverDomain = `http://localhost:8080`
)

var (
	errInvalidGameID error = errors.New(`invalid game id`)
)

type terminalClient struct {
	server *http.Client

	reqChan chan terminalRequest

	me            model.Player
	myCurrentGame model.GameID
	myGames       map[model.GameID]model.Game
}

type termReqType int

const (
	blocking termReqType = iota
	scoreUpdate
	message
	info
	switchGames
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
		reqChan: make(chan terminalRequest, 5),
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

	port, err := findOpenPort()
	if err != nil {
		return err
	}

	tc.startClientServer(&wg, port)
	tc.tellAboutInteraction(&wg, port)
	tc.processUserInput(&wg)

	// Block until forever...?
	wg.Wait()

	return nil
}

func findOpenPort() (int, error) {
	port := 8081
	for ; port < 8180; port++ {
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

		if err != nil {
			continue
		}

		err = ln.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't stop listening on port %q: %s", port, err)
			return 0, err
		}

		break
	}

	return port, nil
}

func (tc *terminalClient) startClientServer(wg *sync.WaitGroup, port int) {
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

		router.POST("/blocking/:gameID", handleBlocking(tc.reqChan))
		router.POST("/message/:gameID", handleMessage(tc.reqChan))
		router.POST("/score/:gameID", handleScoreUpdate(tc.reqChan))

		err = router.Run(fmt.Sprintf("0.0.0.0:%d", port)) // listen and serve on the addr
		fmt.Printf("router.Run error: %+v\n", err)
	}()
}

func (tc *terminalClient) tellAboutInteraction(wg *sync.WaitGroup, port int) {
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
}

func (tc *terminalClient) processUserInput(wg *sync.WaitGroup) {
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

		tc.reqChan <- terminalRequest{
			game: tc.myGames[tc.myCurrentGame],
			msg:  `Starting terminal player`,
			req:  message,
		}
		for req := range tc.reqChan {
			err := tc.processRequest(req)
			if err != nil && err != errInvalidGameID {
				tc.reqChan <- terminalRequest{
					gameID: req.gameID,
					game:   req.game,
					msg:    fmt.Sprintf(`Problem doing action (%s). Try again?`, err.Error()),
					req:    req.req,
				}

			}
		}
	}()
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

	bytes, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		// Keeping this here for debugging
		fmt.Printf("full response: %+v\n%s\n%s\n", response, response.Body, string(bytes))

		contentType := response.Header.Get("Content-Type")
		if strings.Contains(contentType, `text/plain`) {
			return nil, fmt.Errorf("bad response: \"%s\"", string(bytes))
		}

		return nil, fmt.Errorf("bad response from server")
	} else if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (tc *terminalClient) createPlayer() error {
	username, name := tc.getName()
	var reqData = struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	}{username, name}
	b, err := json.Marshal(reqData)
	if err != nil {
		return err
	}
	// TODO do we need makeRequest to set Content-Type to application/json in the header?
	respBytes, err := tc.makeRequest(`POST`, `/create/player`, bytes.NewReader(b))
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

	g, err := jsonutils.UnmarshalGame(respBytes)
	if err != nil {
		return err
	}

	tc.myGames[g.ID] = g

	if tc.myCurrentGame == model.InvalidGameID {
		tc.myCurrentGame = g.ID
		msg := `Joined game with `
		msg += gamePlayerNames(tc.myGames[tc.myCurrentGame])
		msg += `.`
		fmt.Println(msg)
	} else {
		tc.reqChan <- terminalRequest{
			gameID: g.ID,
			req:    switchGames,
		}
	}

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

	tc.reqChan <- terminalRequest{
		msg: fmt.Sprintf(`Knows about %d games`, len(tc.me.Games)),
		req: info,
	}

	for gID := range tc.me.Games {
		g, err := tc.requestGame(gID)
		if err != nil {
			return err
		}

		tc.myGames[gID] = g

		if g.IsOver() {
			playerNames := gamePlayerNames(g)
			gameScore := gameScoreMessage(g, tc.me.ID)
			tc.reqChan <- terminalRequest{
				msg: fmt.Sprintf("Game with %s is over. Final score: \n%s", playerNames, gameScore),
				req: info,
			}
		} else {
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

	g, err := jsonutils.UnmarshalGame(respBytes)
	if err != nil {
		return model.Game{}, err
	}
	tc.myGames[g.ID] = g

	if _, ok := g.BlockingPlayers[tc.me.ID]; ok && g.ID != tc.myCurrentGame {
		tc.reqChan <- terminalRequest{
			gameID: g.ID,
			req:    switchGames,
		}
	}

	return g, nil
}

func (tc *terminalClient) processRequest(req terminalRequest) error {
	switch req.req {
	case info:
		fmt.Println(req.msg)
		return nil
	case switchGames:
		return tc.askToSwitchGames(req.gameID)
	}

	gID := req.gameID
	if req.gameID == model.InvalidGameID {
		gID = req.game.ID
	}
	if gID == model.InvalidGameID {
		fmt.Printf("request does not have valid game ID: %+v\n", req)
		return errInvalidGameID
	}

	switch req.req {
	case blocking:
		fmt.Printf("Blocking message: \"%s\"\n", req.msg)
		err := tc.requestAndSendAction(gID)
		if err != nil {
			return err
		}
	case message:
		fmt.Println(req.msg)
	case scoreUpdate:
		fmt.Println(req.msg)
		tc.printCurrentScore()
	default:
		fmt.Printf("Developer error: req needs a type %+v\n", req)
	}
	return nil
}

func (tc *terminalClient) askToSwitchGames(newGameID model.GameID) error {
	should := true

	msg := `Current game is with `
	msg += gamePlayerNames(tc.myGames[tc.myCurrentGame])
	msg += `. `
	msg += `Do you want to switch to game with `
	msg += gamePlayerNames(tc.myGames[newGameID])
	msg += `? `

	prompt := &survey.Confirm{
		Message: msg,
		Default: true,
	}

	err := survey.AskOne(prompt, &should)
	if err != nil {
		fmt.Printf("survey.AskOne error: %+v\n", err)
		return err
	}

	tc.reqChan <- terminalRequest{
		gameID: newGameID,
		msg:    `Switched to new game`,
		req:    blocking,
	}

	return nil
}

func gamePlayerNames(g model.Game) string {
	msg := ``
	for i, p := range g.Players {
		if i > 0 {
			msg += `, `
		}
		msg += p.Name
	}
	return msg
}
