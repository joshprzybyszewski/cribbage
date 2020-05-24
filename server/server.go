package server

import (
	"context"
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
	db persistence.DB
}

func newCribbageServer(db persistence.DB) *cribbageServer {
	return &cribbageServer{
		db: db,
	}
}

func (cs *cribbageServer) NewRouter() http.Handler {
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
		user.GET("/:username", handleGetUsername)
		user.GET("/:username/game/:gameID", handleGetUsernameGame)
	}

	// Simple group: create
	create := router.Group(`/create`)
	{
		create.POST(`/game`, cs.ginPostCreateGame)
		create.POST(`/player`, cs.ginPostCreatePlayer)
		create.POST(`/interaction`, cs.ginPostCreateInteraction)
	}

	router.GET(`/game/:gameID`, cs.ginGetGame)
	router.GET(`/player/:username`, cs.ginGetPlayer)

	router.POST(`/action`, cs.ginPostAction)

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
	var gameReq network.CreateGameRequest
	err := c.ShouldBindJSON(&gameReq)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	pIDs := make([]model.PlayerID, len(gameReq.PlayerIDs))
	for i, pID := range gameReq.PlayerIDs {
		if pID == model.InvalidPlayerID {
			c.String(http.StatusBadRequest, `Invalid player ID at index %d`, i)
			return
		}
		pIDs[i] = pID
	}

	if len(pIDs) < model.MinPlayerGame || len(pIDs) > model.MaxPlayerGame {
		c.String(http.StatusBadRequest, `Invalid num players: %d`, len(gameReq.PlayerIDs))
		return
	}
	g, err := createGame(context.Background(), cs.db, pIDs)
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
	err = createPlayer(context.Background(), cs.db, player)
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
	pID := cir.PlayerID
	if pID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, `Needs playerId`)
		return
	}

	var pm interaction.PlayerMeans
	switch {
	case len(cir.LocalhostPort) > 0:
		pm = interaction.New(pID, interaction.Means{
			Mode: interaction.Localhost,
			Info: cir.LocalhostPort,
		})
	case len(cir.NPCType) > 0:
		switch cir.NPCType {
		case interaction.Simple, interaction.Calc, interaction.Dumb:
		default:
			c.String(http.StatusBadRequest, `unsupported interaction mode`)
			return
		}
		pm = interaction.New(pID, interaction.Means{
			Mode: interaction.NPC,
			Info: cir.NPCType,
		})
	default:
		c.String(http.StatusBadRequest, `unsupported interaction mode`)
		return
	}

	err = saveInteraction(context.Background(), cs.db, pm)
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
	g, err := getGame(context.Background(), cs.db, gID)
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
	p, err := getPlayer(context.Background(), cs.db, pID)
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
	reqBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}
	action, err := jsonutils.UnmarshalPlayerAction(reqBytes)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	err = handleAction(context.Background(), cs.db, action)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	c.String(http.StatusOK, `action handled`)
}
