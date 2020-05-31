package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

func handleWasmIndex(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		`index.html`,
		gin.H{},
	)
}

func handleWasmGetUser(c *gin.Context) {
	// read the username/displayname from the query params
	// and redirect to /user/:username
	username := c.Query(`username`)

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(`/wasm/user/%s`, username))
}

func (cs *cribbageServer) handleWasmGetUsername(c *gin.Context) {
	ctx := context.Background()
	// serve up a list of games this user is in
	username := c.Param(`username`)
	pID := model.PlayerID(username)

	db, err := cs.dbFactory.New(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `dbFactory.New() error: %s`, err)
		return
	}
	defer db.Close()

	p, err := getPlayer(ctx, db, pID)
	if err != nil {
		c.String(http.StatusInternalServerError, `Error: %s`, err)
		return
	}

	gameNames := make([]struct {
		Desc       string
		MyColor    string
		RedScore   int
		BlueScore  int
		GreenScore int
		GameID     model.GameID
	}, 0, len(p.Games))
	for gID, color := range p.Games {
		g, err := getGame(ctx, db, gID)
		if err != nil {
			if err == persistence.ErrGameNotFound {
				// the player knows about a game that's been deleted
				continue
			}
			c.String(http.StatusInternalServerError, `Error getting game: %s`, err)
			return
		}
		var gameDesc strings.Builder
		gameDesc.WriteString(`Game with`)
		for i, p := range g.Players {
			if i > 0 {
				gameDesc.WriteString(`,`)
			}
			gameDesc.WriteString(` `)
			if p.ID == pID {
				gameDesc.WriteString(`you`)
			} else {
				gameDesc.WriteString(p.Name)
			}
		}
		gameDesc.WriteString(` `)
		gameNames = append(gameNames, struct {
			Desc       string
			MyColor    string
			RedScore   int
			BlueScore  int
			GreenScore int
			GameID     model.GameID
		}{
			Desc:       gameDesc.String(),
			GameID:     gID,
			MyColor:    color.String(),
			RedScore:   g.CurrentScores[model.Red],
			BlueScore:  g.CurrentScores[model.Blue],
			GreenScore: g.CurrentScores[model.Green],
		})
	}

	c.HTML(
		http.StatusOK,
		`user.html`,
		gin.H{
			`displayName`: username,
			`myID`:        string(pID),
			`games`:       gameNames,
		},
	)
}

func (cs *cribbageServer) handleWasmGetUsernameGame(c *gin.Context) { //nolint:gocyclo
	// serve up this game for this user
	pID := model.PlayerID(c.Param(`username`))
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
			c.String(http.StatusBadRequest, `Game (%v) does not exist`, gID)
		}
		c.String(http.StatusInternalServerError, `Problem getting game: %s`, err)
		return
	}

	playerNames := make([]string, 0, len(g.Players))
	nameMap := make(map[model.PlayerID]string, len(g.Players))
	for _, p := range g.Players {
		playerNames = append(playerNames, p.Name)
		nameMap[p.ID] = p.Name
	}

	if _, ok := nameMap[pID]; !ok {
		c.String(http.StatusBadRequest, `Player %v not in game %v`, pID, gID)
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

	cutCard := g.CutCard.String()
	emptyCard := model.Card{}
	if g.CutCard == emptyCard {
		cutCard = ``
	}

	oppHands := []struct {
		Name string
		Hand []struct {
			Card    string
			IsKnown bool
		}
	}{}

	peggedCardMap := make(map[model.Card]struct{}, len(g.PeggedCards))
	peggedCards := make([]string, 0, len(g.PeggedCards))
	for _, pc := range g.PeggedCards {
		peggedCards = append(peggedCards, pc.Card.String())
		peggedCardMap[pc.Card] = struct{}{}
	}

	myHand := make([]struct {
		Card     string
		IsPegged bool
	}, 0, len(g.Hands[pID]))

	for _, c := range g.Hands[pID] {
		_, ok := peggedCardMap[c]
		myHand = append(myHand, struct {
			Card     string
			IsPegged bool
		}{
			Card:     c.String(),
			IsPegged: ok,
		})
	}

	for playerID, hand := range g.Hands {
		if pID == playerID {
			continue
		}
		if len(hand) == 0 {
			continue
		}
		hands := make([]struct {
			Card    string
			IsKnown bool
		}, 0, len(hand))

		for _, c := range hand {
			cStr := ``
			known := false
			if _, ok := peggedCardMap[c]; ok {
				cStr = c.String()
				known = true
			}
			hands = append(hands, struct {
				Card    string
				IsKnown bool
			}{
				Card:    cStr,
				IsKnown: known,
			})
		}

		oppHands = append(oppHands, struct {
			Name string
			Hand []struct {
				Card    string
				IsKnown bool
			}
		}{
			Name: nameMap[playerID],
			Hand: hands,
		})
	}

	cribHand := make([]struct {
		Card    string
		IsKnown bool
	}, 0, len(g.Crib))

	if g.Phase >= model.BuildCribReady {
		for _, c := range g.Crib {
			cStr := ``
			known := false
			if g.Phase >= model.CribCounting {
				cStr = c.String()
				known = true
			}
			cribHand = append(cribHand, struct {
				Card    string
				IsKnown bool
			}{
				Card:    cStr,
				IsKnown: known,
			})
		}
	}

	currentDealerName := nameMap[g.CurrentDealer]
	if g.CurrentDealer == pID {
		currentDealerName = `You`
	}

	c.HTML(
		http.StatusOK,
		`game.html`,
		gin.H{
			`myID`:          string(pID),
			`myColor`:       g.PlayerColors[pID].String(),
			`scores`:        scores,
			`currentDealer`: currentDealerName,
			`myHand`:        myHand,
			`oppHands`:      oppHands,
			`peggedCards`:   peggedCards,
			`currentPeg`:    g.CurrentPeg(),
			`crib`:          cribHand,
			`myCrib`:        g.CurrentDealer == pID,
			`phase`:         g.Phase.String(),
			`cutCard`:       cutCard,
			`playerNames`:   playerNames,
			`game`:          g,
		},
	)
}
