package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

type cribbageServer struct {
	dbService persistence.DB
}

func newCribbageServer(db persistence.DB) *cribbageServer {
	return &cribbageServer{
		dbService: db,
	}
}

func (cs *cribbageServer) NewRouter() http.Handler {
	router := gin.Default()

	// Simple group: create
	create := router.Group(`/create`)
	{
		create.POST(`/game`, cs.ginPostCreateGame)
		create.POST(`/player`, cs.ginPostCreatePlayer)
		create.POST(`/interaction`, cs.ginPostCreateInteraction)
	}

	router.GET(`/game/:gameID`, cs.ginGetGame)
	router.GET(`/player/:username`, cs.ginGetPlayer)

	router.POST(`/action/:gameID`, cs.ginPostAction)

	return router
}

func (cs *cribbageServer) Serve() {
	router := cs.NewRouter()
	eng, ok := router.(*gin.Engine)
	if !ok {
		log.Println(`router type assertion failed`)
	}

	err := eng.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Printf("router.Run errored: %+v\n", err)
	}
}

func (cs *cribbageServer) ginPostCreateGame(c *gin.Context) {
	var cgr network.CreateGameRequest
	err := c.ShouldBindJSON(&cgr)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	pIDs := make([]model.PlayerID, len(cgr.PlayerIDs))
	for i, idStr := range cgr.PlayerIDs {
		pID := model.PlayerID(idStr)
		if pID == model.InvalidPlayerID {
			c.String(http.StatusBadRequest, `Invalid player ID at index %d`, i)
			return
		}
		pIDs[i] = pID
	}

	if len(pIDs) < model.MinPlayerGame || len(pIDs) > model.MaxPlayerGame {
		c.String(http.StatusBadRequest, `Invalid num players: %d`, len(cgr.PlayerIDs))
		return
	}
	g, err := createGame(cs.dbService, pIDs)
	if err != nil {
		c.String(http.StatusInternalServerError, `createGame error: %s`, err)
		return
	}

	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, g)
}

func (cs *cribbageServer) ginPostCreatePlayer(c *gin.Context) {
	var player model.Player
	err := c.ShouldBindJSON(&player)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	if player.ID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, `Username is required`)
		return
	}
	if player.Name == `` {
		c.String(http.StatusBadRequest, `Display name is required`)
		return
	}
	if !model.IsValidPlayerID(player.ID) {
		c.String(http.StatusBadRequest, `Username must be alphanumeric`)
		return
	}
	err = cs.dbService.CreatePlayer(player)
	if err != nil {
		switch err {
		case persistence.ErrPlayerAlreadyExists:
			c.String(http.StatusBadRequest, `Username already exists`)
		default:
			c.String(http.StatusInternalServerError, `Error: %s`, err)
		}
		return
	}
	c.JSON(http.StatusOK, player)
}

func (cs *cribbageServer) ginPostCreateInteraction(c *gin.Context) {
	var cir network.CreateInteractionRequest
	err := c.ShouldBindJSON(&cir)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	pID := model.PlayerID(cir.PlayerID)
	if pID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, `Needs playerId`)
		return
	}

	var mode interaction.Mode
	switch cir.Mode {
	case `localhost`:
		mode = interaction.Localhost
	default:
		c.String(http.StatusBadRequest, `unsupported interaction mode`)
		return
	}

	info := cir.Info
	pm := interaction.New(pID, interaction.Means{
		Mode: mode,
		Info: info,
	})
	err = cs.dbService.SaveInteraction(pm)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	c.String(http.StatusOK, `Updated player interaction`)
}

func (cs *cribbageServer) ginGetGame(c *gin.Context) {
	gID, err := getGameIDFromContext(c)
	if err != nil {
		c.String(http.StatusBadRequest, `Invalid GameID: %v`, err)
		return
	}
	g, err := cs.dbService.GetGame(gID)
	if err != nil {
		if err == persistence.ErrGameNotFound {
			c.String(http.StatusNotFound, `Game not found`)
			return
		}
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, g)
}

func getGameIDFromContext(c *gin.Context) (model.GameID, error) {
	gIDStr := c.Param(`gameID`)
	n, err := strconv.Atoi(gIDStr)
	if err != nil {
		return model.InvalidGameID, err
	}
	return model.GameID(n), nil
}

func (cs *cribbageServer) ginGetPlayer(c *gin.Context) {
	pID := model.PlayerID(c.Param(`username`))
	p, err := cs.dbService.GetPlayer(pID)
	if err != nil {
		if err == persistence.ErrPlayerNotFound {
			c.String(http.StatusNotFound, `Player not found`)
			return
		}
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	// TODO investigate what it'll take to protobuf-ify our models
	c.JSON(http.StatusOK, p)
}

func (cs *cribbageServer) ginPostAction(c *gin.Context) {
	action, err := unmarshalPlayerAction(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	err = handleAction(cs.dbService, action)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	c.String(http.StatusOK, `action handled`)
}

func unmarshalPlayerAction(req *http.Request) (model.PlayerAction, error) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return model.PlayerAction{}, err
	}

	return jsonutils.UnmarshalPlayerAction(reqBody)
}
