package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

type cribbageServer struct {
	db persistence.DB
}

func (cs *cribbageServer) Serve() {
	// TODO add handling to route traffic to the correct method
	// I imagine this will be martini or another REST server router
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Simple group: create
	create := router.Group("/create")
	{
		create.POST("/game", func(c *gin.Context) {
			// TODO ensure we pass in the playerIDs for this game
			var pIDs []model.PlayerID
			g, err := cs.createGame(pIDs)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error: %s", err)
				return
			}
			// TODO investigate what it'll take to protobuf-ify our models
			c.ProtoBuf(http.StatusOK, g)
		})
		create.POST("/player/:name", func(c *gin.Context) {
			name := c.Param("name")
			pID, err := cs.createPlayer(name)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error: %s", err)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"name": name,
				"ID": pID,
			})
		})
		create.POST("/interaction/:playerId", func(c *gin.Context) {
			// TODO ensure this whole thing makes sense...
			pIDStr := c.Param("playerId")
			pID := model.PlayerID(strconv.Atoi(pIDStr))
			err := cs.setInteraction(pID)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error: %s", err)
				return
			}
			c.String(http.StatusOK, "Updated player interaction")
		})
	}

	router.GET("/game/:gameID", func(c *gin.Context) {
		gIDStr := c.Param("gameID")
		gID := model.GameID(strconv.Atoi(gIDStr))
		g, err := cs.GetGame(gID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %s", err)
			return
		}
		// TODO investigate what it'll take to protobuf-ify our models
		c.ProtoBuf(http.StatusOK, g)
	})
	router.GET("/player/:playerID", func(c *gin.Context) {
		pIDStr := c.Param("playerID")
		pID := model.PlayerID(strconv.Atoi(pIDStr))
		p, err := cs.GetPlayer(pID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %s", err)
			return
		}
		// TODO investigate what it'll take to protobuf-ify our models
		c.ProtoBuf(http.StatusOK, p)
	})

	create.POST("/action", func(c *gin.Context) {
		// TODO find out how to pass in the action
		var action model.PlayerAction 
		err := cs.handleAction(action)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error: %s", err)
			return
		}
		// TODO figure out if we need to send back the updated game state
		c.String(http.StatusOK, "action handled")
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

