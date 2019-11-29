package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

type cribbageServer struct {
	db persistence.DB
}

func (cs *cribbageServer) Serve() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Simple group: create
	create := router.Group("/create")
	{
		create.POST("/game/:player1/:player2", cs.ginPostCreateGame)
		create.POST("/game/:player1/:player2/:player3", cs.ginPostCreateGame)
		create.POST("/game/:player1/:player2/:player3/:player4", cs.ginPostCreateGame)
		create.POST("/player/:username/:name", cs.ginPostCreatePlayer)
		create.POST("/interaction/:playerId/:means/:info", cs.ginPostCreateInteraction)
	}

	router.GET("/game/:gameID", cs.ginGetGame)
	router.GET("/player/:username", cs.ginGetPlayer)

	router.POST("/action/:gameID", cs.ginPostAction)

	err := router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Printf("router.Run errored: %+v\n", err)
	}
}

func (cs *cribbageServer) ginPostCreateGame(c *gin.Context) {
	var pIDs []model.PlayerID

	pID := getPlayerID(c, `player1`)
	if pID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, "Needs player1")
		return
	}
	pIDs = append(pIDs, pID)

	pID = getPlayerID(c, `player2`)
	if pID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, "Needs player2")
		return
	}
	pIDs = append(pIDs, pID)

	pID = getPlayerID(c, `player3`)
	if pID != model.InvalidPlayerID {
		pIDs = append(pIDs, pID)
	}

	pID = getPlayerID(c, `player4`)
	if pID != model.InvalidPlayerID {
		pIDs = append(pIDs, pID)
	}

	if len(pIDs) < model.MinPlayerGame || len(pIDs) > model.MaxPlayerGame {
		c.String(http.StatusBadRequest, "Invalid num players: %d", len(pIDs))
		return
	}

	g, err := cs.createGame(pIDs)
	if err != nil {
		c.String(http.StatusInternalServerError, "createGame error: %s", err)
		return
	}

	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, g)
}

func getPlayerID(c *gin.Context, playerParam string) model.PlayerID {
	username, ok := c.Params.Get(playerParam)
	if !ok {
		return model.InvalidPlayerID
	}

	return model.PlayerID(username)
}

func (cs *cribbageServer) ginPostCreatePlayer(c *gin.Context) {
	username := c.Param("username")
	name := c.Param("name")
	player, err := cs.createPlayer(username, name)
	if err != nil {
		switch err {
		case persistence.ErrPlayerAlreadyExists:
			c.String(http.StatusBadRequest, "Username already exists")
		case errInvalidUsername:
			c.String(http.StatusBadRequest, "Username must be alphanumeric")
		default:
			c.String(http.StatusInternalServerError, "Error: %s", err)
		}
		return
	}
	c.JSON(http.StatusOK, player)
}

func (cs *cribbageServer) ginPostCreateInteraction(c *gin.Context) {
	pID := getPlayerID(c, `playerId`)
	if pID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, "Needs playerId")
		return
	}
	means := c.Param(`means`)
	info := c.Param(`info`)
	err := cs.setInteraction(pID, model.InteractionMeans{
		Means: means,
		Info:  info,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: %s", err)
		return
	}
	c.String(http.StatusOK, "Updated player interaction")
}

func (cs *cribbageServer) ginGetGame(c *gin.Context) {
	gIDStr := c.Param("gameID")
	n, err := strconv.Atoi(gIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid GameID: %s", gIDStr)
		return
	}
	gID := model.GameID(n)
	g, err := cs.getGame(gID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: %s", err)
		return
	}
	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, g)
}

func (cs *cribbageServer) ginGetPlayer(c *gin.Context) {
	pID := model.PlayerID(c.Param("username"))
	p, err := cs.getPlayer(pID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: %s", err)
		return
	}
	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, p)
}

func (cs *cribbageServer) ginPostAction(c *gin.Context) {
	action, err := unmarshalPlayerAction(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, "Error: %s", err)
		return
	}

	err = cs.handleAction(action)
	if err != nil {
		c.String(http.StatusBadRequest, "Error: %s", err)
		return
	}

	c.String(http.StatusOK, "action handled")
}

func unmarshalPlayerAction(req *http.Request) (model.PlayerAction, error) {
	// We can store the RawMessage and then switch on the Overcomes type later
	// otherwise Action becomes a map[string]interface{}
	var raw json.RawMessage
	action := model.PlayerAction{
		Action: &raw,
	}
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return model.PlayerAction{}, err
	}
	err = json.Unmarshal(reqBody, &action)
	if err != nil {
		return model.PlayerAction{}, err
	}

	blockerActions := map[model.Blocker]func() interface{}{
		model.DealCards: func() interface{} { return &model.DealAction{} },
		model.CribCard:  func() interface{} { return &model.BuildCribAction{} },
		model.CutCard:   func() interface{} { return &model.CutDeckAction{} },
		model.PegCard:   func() interface{} { return &model.PegAction{} },
		model.CountHand: func() interface{} { return &model.CountHandAction{} },
		model.CountCrib: func() interface{} { return &model.CountCribAction{} },
	}

	subActionFn, ok := blockerActions[action.Overcomes]
	if !ok {
		return model.PlayerAction{}, errors.New(`unknown action type`)
	}
	subAction := subActionFn()

	if err := json.Unmarshal(raw, subAction); err != nil {
		return model.PlayerAction{}, err
	}

	switch t := subAction.(type) {
	case *model.DealAction:
		action.Action = *t
	case *model.BuildCribAction:
		action.Action = *t
	case *model.CutDeckAction:
		action.Action = *t
	case *model.PegAction:
		action.Action = *t
	case *model.CountHandAction:
		action.Action = *t
	case *model.CountCribAction:
		action.Action = *t
	}

	return action, nil
}
