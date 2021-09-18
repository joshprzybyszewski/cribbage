package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/logic/suggestions"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

type cribbageServer struct {
	dbFactory persistence.DBFactory
}

func newCribbageServer(dbFactory persistence.DBFactory) *cribbageServer {
	return &cribbageServer{
		dbFactory: dbFactory,
	}
}

func (cs *cribbageServer) NewRouter() http.Handler {
	router := gin.Default()

	// health check route
	router.GET(`/health`, func(c *gin.Context) {
		c.String(http.StatusOK, `Healthy!`)
	})

	// Simple group: create
	create := router.Group(`/create`)
	{
		create.POST(`/game`, cs.ginPostCreateGame)
		create.POST(`/player`, cs.ginPostCreatePlayer)
		create.POST(`/interaction`, cs.ginPostCreateInteraction)
	}

	router.GET(`/game/:gameID`, cs.ginGetGame)

	// Simple group: games
	game := router.Group(`/games`)
	{
		game.GET(`/active`, cs.ginGetActiveGamesForPlayer)
	}

	// Simple group: player
	player := router.Group(`/player`)
	{
		player.GET(`/:username`, cs.ginGetPlayer)
	}

	router.POST(`/action`, cs.ginPostAction)

	// Simple group: suggest
	suggest := router.Group(`/suggest`)
	{
		suggest.GET(`/hand`, cs.ginGetSuggestHand)
	}

	return router
}

func (cs *cribbageServer) addWasmHandlers(router *gin.Engine) {
	router.LoadHTMLGlob(`templates/*`)
	router.Static(`/assets`, `./assets`)

	wasm := router.Group(`/wasm`)
	{
		wasm.GET(`/`, handleWasmIndex)

		// Simple group: user. Used for serving pages affiliated with a given user
		user := wasm.Group(`/user`)
		{
			user.GET(`/`, handleWasmGetUser)
			user.GET(`/:username`, cs.handleWasmGetUsername)
			user.GET(`/:username/game/:gameID`, cs.handleWasmGetUsernameGame)
		}
	}
}

func (cs *cribbageServer) addReactHandlers(router *gin.Engine) {
	// Serve frontend React static files
	router.Use(static.Serve(`/`, static.LocalFile(`./client/build`, true)))
}

type serveConfig struct {
	includeStaticResources bool
}

func (cs *cribbageServer) Serve(config serveConfig) http.Handler {
	router := cs.NewRouter()
	if config.includeStaticResources {
		eng, ok := router.(*gin.Engine)
		if !ok {
			log.Println(`router type assertion failed`)
		}
		cs.addWasmHandlers(eng)
		cs.addReactHandlers(eng)
	}
	return router
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

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	g, err := createGame(ctx, db, pIDs)
	if err != nil {
		c.String(http.StatusInternalServerError, `createGame error: %s`, err)
		return
	}

	c.JSON(http.StatusOK, network.ConvertToCreateGameResponse(g))
}

// POST /create/player
func (cs *cribbageServer) ginPostCreatePlayer(c *gin.Context) {
	var cpr network.CreatePlayerRequest
	err := c.ShouldBindJSON(&cpr)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	if cpr.Player.ID == model.InvalidPlayerID {
		c.String(http.StatusBadRequest, `Username is required`)
		return
	}
	if cpr.Player.Name == `` {
		c.String(http.StatusBadRequest, `Display name is required`)
		return
	}
	if !model.IsValidPlayerID(cpr.Player.ID) {
		c.String(http.StatusBadRequest, `Username must be alphanumeric`)
		return
	}

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	p := model.Player{
		ID:   cpr.Player.ID,
		Name: cpr.Player.Name,
	}
	err = createPlayer(ctx, db, p)
	if err != nil {
		switch err {
		case persistence.ErrPlayerAlreadyExists:
			c.String(http.StatusBadRequest, `Username already exists`)
		default:
			c.String(http.StatusInternalServerError, `Error: %s`, err)
		}
		return
	}
	c.JSON(http.StatusOK, network.ConvertToCreatePlayerResponse(p))
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

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	err = saveInteraction(ctx, db, pm)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	c.String(http.StatusOK, `Updated player interaction`)
}

// GET /game/:gameID?player=<playerID>
func (cs *cribbageServer) ginGetGame(c *gin.Context) {
	gID, err := getGameIDFromContext(c)
	if err != nil {
		c.String(http.StatusBadRequest, `Invalid GameID: %v`, err)
		return
	}

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	g, err := getGame(ctx, db, gID)
	if err != nil {
		if err == persistence.ErrGameNotFound {
			c.String(http.StatusNotFound, `Game not found`)
			return
		}
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}

	pID := c.Query(`player`)
	if pID == `` {
		resp := network.ConvertToGetGameResponse(g)
		c.JSON(http.StatusOK, resp)
		return
	}
	resp, err := network.ConvertToGetGameResponseForPlayer(g, model.PlayerID(pID))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func getGameIDFromContext(c *gin.Context) (model.GameID, error) {
	gIDStr := c.Param(`gameID`)
	n, err := strconv.Atoi(gIDStr)
	if err != nil {
		return model.InvalidGameID, err
	}
	return model.GameID(n), nil
}

// GET /player/:username
func (cs *cribbageServer) ginGetPlayer(c *gin.Context) {
	pID := model.PlayerID(c.Param(`username`))

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	p, err := getPlayer(ctx, db, pID)
	if err != nil {
		if err == persistence.ErrPlayerNotFound {
			c.String(http.StatusNotFound, `Player not found`)
			return
		}
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	resp := network.ConvertToGetPlayerResponse(p)
	c.JSON(http.StatusOK, resp)
}

// GET /games/active?playerID=pID
func (cs *cribbageServer) ginGetActiveGamesForPlayer(c *gin.Context) {
	pID := model.PlayerID(c.Query(`playerID`))
	if len(pID) == 0 {
		c.String(http.StatusBadRequest, `Requires playerID`)
		return
	}

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	p, err := getPlayer(ctx, db, pID)
	if err != nil {
		if err == persistence.ErrPlayerNotFound {
			c.String(http.StatusNotFound, `Player not found`)
			return
		}
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}
	games := make(map[model.GameID]model.Game, len(p.Games))
	for gID := range p.Games {
		mg, err := getGame(ctx, db, gID)
		if err != nil {
			c.String(http.StatusInternalServerError, `Error: %s`, err)
			return
		} else if mg.IsOver() {
			continue
		}
		games[gID] = mg
	}
	resp := network.ConvertToGetActiveGamesForPlayerResponse(p, games)
	c.JSON(http.StatusOK, resp)
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

	ctx := context.Background()
	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	err = handleAction(ctx, db, action)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	c.String(http.StatusOK, `action handled`)
}

// GET /suggest/hand?dealt=<cards>
func (cs *cribbageServer) ginGetSuggestHand(c *gin.Context) {

	hand, err := convertToHand(c.Query(`dealt`))
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	summaries, err := suggestions.GetAllTosses(hand)
	if err != nil {
		c.String(http.StatusBadRequest, `Error: %s`, err)
		return
	}

	resp := network.ConvertToGetSuggestHandResponse(summaries)
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].HandPts.Avg > resp[j].HandPts.Avg
	})
	c.JSON(http.StatusOK, resp)
}

func convertToHand(input interface{}) ([]model.Card, error) {
	inputStr, ok := input.(string)
	if !ok || inputStr == `` {
		return nil, errors.New(`empty dealt hand`)
	}
	var cards []model.Card

	cardStrs := strings.Split(inputStr, `,`)
	for _, cs := range cardStrs {
		c, err := model.NewCardFromExternalString(cs)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	return cards, nil
}
