package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

func handleIndex(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Home Page",
		},
	)
}

func handleGetUser(c *gin.Context) {
	// read the username/displayname from the query params
	// and redirect to /user/:username
	username := c.Query(`username`)
	name := c.Query(`displayName`)

	fmt.Printf("username: %s, displayname: %s\n", username, name)

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(`/user/%s`, username))
}

func handlePostUser(c *gin.Context) {
	// read the username/displayname from the request body
	// create the user
	// and redirect to /user/:username
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Must send appropriate POST data")
		return
	}
	var username, name string

	fmt.Printf("username: %s, displayname: %s\n", username, name)
	fmt.Printf("reqBody: %s\n", reqBody)

	player, err := createPlayerFromNames(username, name)
	if err != nil && err != persistence.ErrPlayerAlreadyExists {
		switch err {
		case errInvalidUsername:
			c.String(http.StatusBadRequest, "Username must be alphanumeric")
		default:
			c.String(http.StatusInternalServerError, "Error: %s", err)
		}
		return
	}

	c.Redirect(http.StatusOK, fmt.Sprintf(`/user/%s`, player.ID))
}

type userGame struct {
	Desc       string
	MyColor    string
	RedScore   int
	BlueScore  int
	GreenScore int
	GameID     model.GameID
}

func handleGetUsername(c *gin.Context) {
	// serve up a list of games this user is in
	username := c.Param("username")
	pID := model.PlayerID(username)
	p, err := getPlayer(pID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: %s", err)
		return
	}

	gameNames := make([]userGame, 0, len(p.Games))
	for gID, color := range p.Games {
		g, err := getGame(gID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error getting game: %s", err)
			return
		}
		gameDesc := `Against `
		for _, p := range g.Players {
			if p.ID == pID {
				continue
			}
			// TODO support three and four player games
			gameDesc += p.Name
		}
		gameDesc += ` `
		gameNames = append(gameNames, userGame{
			Desc:       gameDesc,
			GameID:     gID,
			MyColor:    color.String(),
			RedScore:   g.CurrentScores[model.Red],
			BlueScore:  g.CurrentScores[model.Blue],
			GreenScore: g.CurrentScores[model.Green],
		})
	}

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"displayName": username,
			"myID":        string(pID),
			"games":       gameNames,
		},
	)
}

func handleGetUsernameGame(c *gin.Context) {
	// serve up this game for this user
	pID := model.PlayerID(c.Param("username"))
	gIDStr := c.Param("gameID")
	n, err := strconv.Atoi(gIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid GameID: %s", gIDStr)
		return
	}
	gID := model.GameID(n)

	g, err := getGame(gID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Problem getting game: %s", err)
		return
	}

	scores := []struct {
		Color string
		Score int
	}{}

	for color, score := range g.CurrentScores {
		scores = append(scores, struct {
			Color string
			Score int
		}{
			Color: color.String(),
			Score: score,
		})
	}

	myHand := []struct {
		Card string
	}{}

	for _, c := range g.Hands[pID] {
		myHand = append(myHand, struct {
			Card string
		}{
			Card: c.String(),
		})
	}

	cutCard := g.CutCard.String()
	emptyCard := model.Card{}
	if g.CutCard == emptyCard {
		cutCard = `not cut`
	}

	c.HTML(
		http.StatusOK,
		"game.html",
		gin.H{
			"myID":    string(pID),
			"scores":  scores,
			"myHand":  myHand,
			"phase":   g.Phase.String(),
			"cutCard": cutCard,
			"game":    g,
		},
	)
}
