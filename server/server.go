package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

type cribbageServer struct {
}

func (cs *cribbageServer) Serve() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", handleIndex)

	// Simple group: user. Used for serving pages affiliated with a given user
	user := router.Group("/user")
	{
		user.GET("/", handleGetUser)
		user.POST("/", handlePostUser)
		user.GET("/:username", handleGetUsername)
		user.GET("/:username/game/:gameID", handleGetUsernameGame)
	}

	// Simple group: create
	create := router.Group("/create")
	{
		create.POST("/game/:player1/:player2", cs.ginPostCreateGame)
		create.POST("/game/:player1/:player2/:player3", cs.ginPostCreateGame)
		create.POST("/game/:player1/:player2/:player3/:player4", cs.ginPostCreateGame)
		create.POST("/player/:username/:name", cs.ginPostCreatePlayer)
		create.POST("/interaction/:playerId/:mode/:info", cs.ginPostCreateInteraction)
	}

	router.GET("/game/:gameID", cs.ginGetGame)
	router.GET("/player/:username", cs.ginGetPlayer)

	router.POST("/action/:gameID", cs.ginPostAction)

	err := router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Printf("router.Run errored: %+v\n", err)
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

	g, err := createGameFromIDs(pIDs)
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
	player, err := createPlayerFromNames(username, name)
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

	var mode interaction.Mode
	switch c.Param(`mode`) {
	case `localhost`:
		mode = interaction.Localhost
	default:
		c.String(http.StatusBadRequest, "unsupported interaction mode")
		return
	}

	info := c.Param(`info`)
	err := setInteraction(pID, interaction.Means{
		Mode: mode,
		Info: info,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: %s", err)
		return
	}
	c.String(http.StatusOK, "Updated player interaction")
}

func (cs *cribbageServer) ginGetGame(c *gin.Context) {
	gID, err := getGameIDFromContext(c)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid GameID: %v", err)
		return
	}
	g, err := getGame(gID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: %s", err)
		return
	}
	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, g)
}

func getGameIDFromContext(c *gin.Context) (model.GameID, error) {
	gIDStr := c.Param("gameID")
	n, err := strconv.Atoi(gIDStr)
	if err != nil {
		return model.InvalidGameID, err
	}
	return model.GameID(n), nil
}

func (cs *cribbageServer) ginGetPlayer(c *gin.Context) {
	pID := model.PlayerID(c.Param("username"))
	p, err := getPlayer(pID)
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

	err = handlePlayerAction(action)
	if err != nil {
		c.String(http.StatusBadRequest, "Error: %s", err)
		return
	}

	c.String(http.StatusOK, "action handled")
}

func unmarshalPlayerAction(req *http.Request) (model.PlayerAction, error) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return model.PlayerAction{}, err
	}

	return jsonutils.UnmarshalPlayerAction(reqBody)
}
