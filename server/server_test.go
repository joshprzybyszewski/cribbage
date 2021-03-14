package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
)

func performRequest(r http.Handler, method, path string, body io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, nil
}

func readBody(t *testing.T, r io.Reader, v interface{}) {
	bs, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	err = json.Unmarshal(bs, v)
	require.NoError(t, err)
}

func readError(t *testing.T, w *httptest.ResponseRecorder) string {
	errMsgBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	return string(errMsgBytes)
}

func prepareBody(t *testing.T, v interface{}) io.Reader {
	reqBytes, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(reqBytes)
}

func newServerAndRouter(_ *testing.T) (*cribbageServer, http.Handler) {
	// first make sure the db is completely cleared
	dbf := memory.NewFactory()
	memory.Clear()
	cs := newCribbageServer(dbf)
	router := cs.NewRouter()
	return cs, router
}

func seedPlayers(t *testing.T, dbf persistence.DBFactory, n int) []model.PlayerID {
	db, err := dbf.New(context.Background())
	require.NoError(t, err)
	defer db.Close()
	pIDs := make([]model.PlayerID, n)
	for i := range pIDs {
		idStr := fmt.Sprintf(`p%d`, i+1)
		err := db.CreatePlayer(model.Player{
			ID:   model.PlayerID(idStr),
			Name: `name`,
		})
		require.NoError(t, err)
		pIDs[i] = model.PlayerID(idStr)
	}
	return pIDs
}

func TestGinPostCreatePlayer(t *testing.T) {
	type testRequest struct {
		req     network.CreatePlayerRequest
		expCode int
		expErr  string
	}
	testCases := []struct {
		msg  string
		reqs []testRequest
	}{{
		msg: `normal stuff`,
		reqs: []testRequest{{
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   `abc`,
					Name: `def`,
				},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}},
	}, {
		msg: `username with weird characters shouldn't return 404`,
		reqs: []testRequest{{
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   `#`,
					Name: `#`,
				},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username must be alphanumeric`,
		}},
	}, {
		msg: `creating the same player errors`,
		reqs: []testRequest{{
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   `abc`,
					Name: `def`,
				},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}, {
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   `abc`,
					Name: `def`,
				},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username already exists`,
		}},
	}, {
		msg: `empty username`,
		reqs: []testRequest{{
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   ``,
					Name: `def`,
				},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username is required`,
		}},
	}, {
		msg: `empty display name`,
		reqs: []testRequest{{
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   `abc`,
					Name: ``,
				},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Display name is required`,
		}},
	}, {
		msg: `send wrong JSON data - this is equivalent to PlayerID and DispName being empty`,
		reqs: []testRequest{{
			req: network.CreatePlayerRequest{
				Player: network.Player{
					ID:   ``,
					Name: ``,
				},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username is required`,
		}},
	}}
	for _, tc := range testCases {
		_, router := newServerAndRouter(t)

		// make the requests
		for _, r := range tc.reqs {
			body := prepareBody(t, r.req)
			w, err := performRequest(router, `POST`, `/create/player`, body)
			require.NoError(t, err)
			// verify
			require.Equal(t, r.expCode, w.Code)
			if r.expCode != http.StatusOK {
				errMsg := readError(t, w)
				assert.Equal(t, r.expErr, errMsg)
				continue
			}
			expPlayer := model.Player{
				ID:   r.req.Player.ID,
				Name: r.req.Player.Name,
			}
			var playerResp network.CreatePlayerResponse
			readBody(t, w.Body, &playerResp)
			player := model.Player{
				ID:   playerResp.Player.ID,
				Name: playerResp.Player.Name,
			}
			assert.NoError(t, err)
			assert.Equal(t, expPlayer, player, tc.msg)
		}
	}
}
func TestGinPostCreateGame(t *testing.T) {
	testCases := []struct {
		msg     string
		pIDs    []string
		expCode int
		expErr  string
	}{{
		msg:     `two player game`,
		pIDs:    []string{`p1`, `p2`},
		expCode: http.StatusOK,
		expErr:  ``,
	}, {
		msg:     `three player game`,
		pIDs:    []string{`p1`, `p2`, `p3`},
		expCode: http.StatusOK,
		expErr:  ``,
	}, {
		msg:     `four player game`,
		pIDs:    []string{`p1`, `p2`, `p3`, `p4`},
		expCode: http.StatusOK,
		expErr:  ``,
	}, {
		msg:     `one player game is an error`,
		pIDs:    []string{`p1`},
		expCode: http.StatusBadRequest,
		expErr:  `Invalid num players: 1`,
	}, {
		msg:     `five player game`,
		pIDs:    []string{`p1`, `p2`, `p3`, `p4`, `p5`},
		expCode: http.StatusBadRequest,
		expErr:  `Invalid num players: 5`,
	}, {
		msg:     `zero player game`,
		pIDs:    []string{},
		expCode: http.StatusBadRequest,
		expErr:  `Invalid num players: 0`,
	}, {
		msg:     `missing player id`,
		pIDs:    []string{`p1`, `p2`, ``, `p4`},
		expCode: http.StatusBadRequest,
		expErr:  `Invalid player ID at index 2`,
	}, {
		msg:     `invalid player id`,
		pIDs:    []string{`p1`, `p2`, `#`, `p4`},
		expCode: http.StatusInternalServerError,
		expErr:  `createGame error: player not found`,
	}, {
		msg:     `create a game with nonexistent players`,
		pIDs:    []string{`p1`, `p6`},
		expCode: http.StatusInternalServerError,
		expErr:  `createGame error: player not found`,
	}, {
		msg:     `bad request body - equivalent to having no player IDs`,
		pIDs:    []string{},
		expCode: http.StatusBadRequest,
		expErr:  `Invalid num players: 0`,
	}}
	cs, router := newServerAndRouter(t)
	// seed the db with players
	seedPlayers(t, cs.dbFactory, 5)
	for _, tc := range testCases {
		cgr := network.CreateGameRequest{}
		cgr.PlayerIDs = make([]model.PlayerID, len(tc.pIDs))
		for i, id := range tc.pIDs {
			cgr.PlayerIDs[i] = model.PlayerID(id)
		}
		// make the request
		body := prepareBody(t, cgr)
		w, err := performRequest(router, `POST`, `/create/game`, body)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.expCode, w.Code)
		if tc.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.expErr, errMsg)
			continue
		}
		var gameResp network.CreateGameResponse
		readBody(t, w.Body, &gameResp)
		// verify the players are in the game
		require.Len(t, gameResp.Players, len(cgr.PlayerIDs))
		for _, p := range gameResp.Players {
			assert.Contains(t, cgr.PlayerIDs, p.ID)
		}
	}
}
func TestGinPostCreateInteraction(t *testing.T) {
	testCases := []struct {
		msg     string
		pIDs    []string
		reqData network.CreateInteractionRequest
		expCode int
		expErr  string
	}{{
		msg: `missing player ID`,
		reqData: network.CreateInteractionRequest{
			PlayerID:      ``,
			LocalhostPort: `1234`,
		},
		expCode: http.StatusBadRequest,
		expErr:  `Needs playerId`,
	}, {
		msg:     `bad request body - equivalent to an empty network.CreateInteractionRequest`,
		reqData: network.CreateInteractionRequest{},
		expCode: http.StatusBadRequest,
		expErr:  `Needs playerId`,
	}, {
		msg: `good request`,
		reqData: network.CreateInteractionRequest{
			PlayerID:      `p1`,
			LocalhostPort: `1234`,
		},
		expCode: http.StatusOK,
		expErr:  ``,
	}, {
		msg: `unsupported interaction mode`,
		reqData: network.CreateInteractionRequest{
			PlayerID: `p1`,
		},
		expCode: http.StatusBadRequest,
		expErr:  `unsupported interaction mode`,
	}}
	cs, router := newServerAndRouter(t)
	seedPlayers(t, cs.dbFactory, 5)
	for _, tc := range testCases {
		// make the request
		body := prepareBody(t, tc.reqData)
		w, err := performRequest(router, `POST`, `/create/interaction`, body)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.expCode, w.Code)
		if tc.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.expErr, errMsg)
			continue
		}
		bs, err := ioutil.ReadAll(w.Body)
		require.NoError(t, err)
		msg := string(bs)
		// verify the players are in the game
		assert.Equal(t, `Updated player interaction`, msg)
	}
}
func TestGinGetGame(t *testing.T) {
	createTestGame := func(t *testing.T, cs *cribbageServer, pIDs []model.PlayerID) model.Game {
		ctx := context.Background()
		db, err := cs.dbFactory.New(ctx)
		require.NoError(t, err)
		defer db.Close()
		g, err := createGame(ctx, db, pIDs)
		require.NoError(t, err)
		return g
	}

	testCases := []struct {
		msg     string
		setup   func(cs *cribbageServer, pIDs []model.PlayerID) (model.Game, string)
		gameID  string
		expCode int
		expErr  string
	}{{
		msg: `bad game ID`,
		setup: func(cs *cribbageServer, pIDs []model.PlayerID) (model.Game, string) {
			g := createTestGame(t, cs, pIDs)
			return g, `/game/123zzz`
		},
		expCode: http.StatusBadRequest,
		expErr:  `Invalid GameID: strconv.Atoi: parsing "123zzz": invalid syntax`,
	}, {
		msg: `normal request`,
		setup: func(cs *cribbageServer, pIDs []model.PlayerID) (model.Game, string) {
			g := createTestGame(t, cs, pIDs)
			return g, fmt.Sprintf(`/game/%d`, g.ID)
		},
		expCode: http.StatusOK,
		expErr:  ``,
	}, {
		msg: `nonexistent game`,
		setup: func(cs *cribbageServer, pIDs []model.PlayerID) (model.Game, string) {
			g := createTestGame(t, cs, pIDs)
			return g, `/game/123`
		},
		expCode: http.StatusNotFound,
		expErr:  `Game not found`,
	}}
	cs, router := newServerAndRouter(t)
	pIDs := seedPlayers(t, cs.dbFactory, 2)
	for _, tc := range testCases {
		g, url := tc.setup(cs, pIDs)
		w, err := performRequest(router, `GET`, url, nil)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.expCode, w.Code)
		if tc.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.expErr, errMsg)
			continue
		}
		var gameResp network.GetGameResponse
		readBody(t, w.Body, &gameResp)
		assert.Equal(t, g.ID, gameResp.ID)
	}
}
func TestGinGetPlayer(t *testing.T) {
	testCases := []struct {
		msg      string
		playerID string
		expCode  int
		expErr   string
	}{{
		msg:      `good request`,
		playerID: `p1`,
		expCode:  http.StatusOK,
		expErr:   ``,
	}, {
		msg:      `nonexistent player`,
		playerID: `p9`,
		expCode:  http.StatusNotFound,
		expErr:   `Player not found`,
	}}
	cs, router := newServerAndRouter(t)
	seedPlayers(t, cs.dbFactory, 2)
	for _, tc := range testCases {
		// make the request
		url := `/player/` + tc.playerID
		w, err := performRequest(router, `GET`, url, nil)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.expCode, w.Code)
		if tc.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.expErr, errMsg)
			continue
		}
		var player network.GetPlayerResponse
		readBody(t, w.Body, &player)
		assert.Equal(t, model.PlayerID(tc.playerID), player.Player.ID)
	}
}

func TestGinPostAction(t *testing.T) {
	type request struct {
		action  model.PlayerAction
		expCode int
		expErr  string
	}
	testCases := []struct {
		msg  string
		reqs []request
	}{{
		msg: `invalid action type`,
		reqs: []request{{
			action: model.PlayerAction{
				ID:        `p1`,
				Overcomes: 123,
				Action:    ``,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Error: unknown action type`,
		}},
	}, {
		msg: `try to do an action at an inappropriate time`,
		reqs: []request{{
			action: model.PlayerAction{
				ID:        `p1`,
				Overcomes: model.CountCrib,
				Action: model.CountCribAction{
					Pts: 2,
				},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Error: Should overcome DealCards, but overcomes CountCrib`,
		}},
	}, {
		msg: `play a few actions`,
		reqs: []request{{
			action: model.PlayerAction{
				ID:        `p1`,
				Overcomes: model.DealCards,
				Action: model.DealAction{
					NumShuffles: 1,
				},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}, {
			// note: these BuildCribActions are attempting to put cards into the crib that
			// these players likely don't have, but we're just testing the router here
			action: model.PlayerAction{
				ID:        `p1`,
				Overcomes: model.CribCard,
				Action: model.BuildCribAction{
					Cards: []model.Card{{
						Suit:  model.Hearts,
						Value: 1,
					}, {
						Suit:  model.Spades,
						Value: 3,
					}},
				},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}, {
			action: model.PlayerAction{
				ID:        `p2`,
				Overcomes: model.CribCard,
				Action: model.BuildCribAction{
					Cards: []model.Card{{
						Suit:  model.Hearts,
						Value: 1,
					}, {
						Suit:  model.Spades,
						Value: 3,
					}},
				},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}},
	}}
	cs, router := newServerAndRouter(t)
	pIDs := seedPlayers(t, cs.dbFactory, 2)
	for _, tc := range testCases {
		// create a game
		ctx := context.Background()
		db, err := cs.dbFactory.New(ctx)
		require.NoError(t, err)
		defer db.Close()

		game, err := createGame(ctx, db, pIDs)
		require.NoError(t, err)

		actionsCompleted := 0
		for _, r := range tc.reqs {
			r.action.GameID = game.ID
			// make the request
			body := prepareBody(t, r.action)
			w, err := performRequest(router, `POST`, `/action`, body)
			require.NoError(t, err)
			// verify
			require.Equal(t, r.expCode, w.Code)
			if r.expCode != http.StatusOK {
				errMsg := readError(t, w)
				assert.Equal(t, r.expErr, errMsg)
				continue
			}
			actionsCompleted++
			bs, err := ioutil.ReadAll(w.Body)
			require.NoError(t, err)
			msg := string(bs)
			// verify the players are in the game
			assert.Equal(t, `action handled`, msg)
			g, err := db.GetGame(game.ID)
			require.NoError(t, err)
			assert.Equal(t, actionsCompleted, len(g.Actions))
		}
	}
}

func TestGinGetSuggestHand(t *testing.T) {
	testCases := []struct {
		msg      string
		url      string
		expCode  int
		expErr   string
		expSuggs []network.GetSuggestHandResponse
	}{{
		msg:     `nothing dealt`,
		url:     `/suggest/hand?dealt=`,
		expCode: http.StatusBadRequest,
		expErr:  `Error: empty dealt hand`,
	}, {
		msg:     `good request`,
		url:     `/suggest/hand?dealt=JH,KH,QH,9H,10H`,
		expCode: http.StatusOK,
		expErr:  ``,
		expSuggs: []network.GetSuggestHandResponse{{
			Hand: []string{`JH`, `QH`, `9H`, `10H`},
			Toss: []string{`KH`},
			HandPts: network.PointStats{
				Min:    8,
				Avg:    10.702127659574469,
				Median: 10,
				Max:    16,
			},
			CribPts: network.PointStats{
				Min:    0,
				Avg:    4.1999159027836175,
				Median: 4,
				Max:    28,
			},
		}, {
			Hand: []string{`JH`, `KH`, `QH`, `10H`},
			Toss: []string{`9H`},
			HandPts: network.PointStats{
				Min:    8,
				Avg:    10.617021276595745,
				Median: 9,
				Max:    18,
			},
			CribPts: network.PointStats{
				Min:    0,
				Avg:    4.69022510021585,
				Median: 4,
				Max:    24,
			},
		}, {
			Hand: []string{`JH`, `KH`, `9H`, `10H`},
			Toss: []string{`QH`},
			HandPts: network.PointStats{
				Min:    7,
				Avg:    9.319148936170214,
				Median: 9,
				Max:    15,
			},
			CribPts: network.PointStats{
				Min:    0,
				Avg:    4.337364393238584,
				Median: 4,
				Max:    28,
			},
		}, {
			Hand: []string{`JH`, `KH`, `QH`, `9H`},
			Toss: []string{`10H`},
			HandPts: network.PointStats{
				Min:    7,
				Avg:    9.23404255319149,
				Median: 9,
				Max:    15,
			},
			CribPts: network.PointStats{
				Min:    0,
				Avg:    4.5008830207720125,
				Median: 4,
				Max:    28,
			},
		}, {
			Hand: []string{`KH`, `QH`, `9H`, `10H`},
			Toss: []string{`JH`},
			HandPts: network.PointStats{
				Min:    4,
				Avg:    5.9361702127659575,
				Median: 6,
				Max:    11,
			},
			CribPts: network.PointStats{
				Min:    0,
				Avg:    4.641493566562946,
				Median: 4,
				Max:    29,
			},
		}},
	}}
	_, router := newServerAndRouter(t)
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) {
			// make the request
			w, err := performRequest(router, `GET`, tc.url, nil)
			require.NoError(t, err)
			// verify
			require.Equal(t, tc.expCode, w.Code)
			if tc.expCode != http.StatusOK {
				errMsg := readError(t, w)
				assert.Equal(t, tc.expErr, errMsg)
				return
			}

			var suggs []network.GetSuggestHandResponse
			readBody(t, w.Body, &suggs)

			assert.Equal(t, tc.expSuggs, suggs)
		})
	}
}
