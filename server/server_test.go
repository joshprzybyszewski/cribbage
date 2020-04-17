package server

import (
	"bytes"
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

func newServerAndRouter() (*cribbageServer, http.Handler) {
	// first make sure the db is completely cleared
	db := memory.New()
	memory.Clear()
	cs := newCribbageServer(db)
	router := cs.NewRouter()
	return cs, router
}

func seedPlayers(t *testing.T, db persistence.DB, n int) []model.PlayerID {
	pIDs := make([]model.PlayerID, n)
	for i := range pIDs {
		idStr := fmt.Sprintf(`p%d`, i+1)
		_, err := createPlayer(db, idStr, `name`)
		require.NoError(t, err)
		pIDs[i] = model.PlayerID(idStr)
	}
	return pIDs
}

type testRequest struct {
	// body is raw json
	reqData interface{}
	body    string
	expCode int
	expErr  string
}

func TestGinPostCreatePlayer(t *testing.T) {
	type testRequest struct {
		PlayerID string `json:"id"`
		DispName string `json:"n"`
		expCode  int
		expErr   string
	}
	testCases := []struct {
		msg  string
		reqs []testRequest
	}{{
		msg: `normal stuff`,
		reqs: []testRequest{{
			PlayerID: `abc`,
			DispName: `def`,
			expCode:  http.StatusOK,
			expErr:   ``,
		}},
	}, {
		msg: `username with weird characters shouldn't return 404`,
		reqs: []testRequest{{
			PlayerID: `#`,
			DispName: `#`,
			expCode:  http.StatusBadRequest,
			expErr:   `Username must be alphanumeric`,
		}},
	}, {
		msg: `creating the same player errors`,
		reqs: []testRequest{{
			PlayerID: `abc`,
			DispName: `def`,
			expCode:  http.StatusOK,
			expErr:   ``,
		}, {
			PlayerID: `abc`,
			DispName: `def`,
			expCode:  http.StatusBadRequest,
			expErr:   `Username already exists`,
		}},
	}, {
		msg: `empty username`,
		reqs: []testRequest{{
			PlayerID: ``,
			DispName: `def`,
			expCode:  http.StatusBadRequest,
			expErr:   `Username is required`,
		}},
	}, {
		msg: `empty display name`,
		reqs: []testRequest{{
			PlayerID: `abc`,
			DispName: ``,
			expCode:  http.StatusBadRequest,
			expErr:   `Display name is required`,
		}},
	}, {
		msg: `send wrong JSON data - this is equivalent to PlayerID and DispName being empty`,
		reqs: []testRequest{{
			PlayerID: ``,
			DispName: ``,
			expCode:  http.StatusBadRequest,
			expErr:   `Username is required`,
		}},
	}}
	for _, tc := range testCases {
		_, router := newServerAndRouter()

		// make the requests
		for _, r := range tc.reqs {
			body := prepareBody(t, r)
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
				ID:   model.PlayerID(r.PlayerID),
				Name: r.DispName,
			}
			var player model.Player
			readBody(t, w.Body, &player)
			assert.NoError(t, err)
			assert.Equal(t, expPlayer, player)
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
	cs, router := newServerAndRouter()
	// seed the db with players
	for i := 0; i < 5; i++ {
		_, err := createPlayer(cs.dbService, fmt.Sprintf(`p%d`, i+1), `name`)
		require.NoError(t, err)
	}
	for _, tc := range testCases {
		g := model.Game{}
		g.Players = make([]model.Player, len(tc.pIDs))
		for i, id := range tc.pIDs {
			g.Players[i] = model.Player{ID: model.PlayerID(id)}
		}
		// make the request
		body := prepareBody(t, g)
		w, err := performRequest(router, `POST`, `/create/game`, body)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.expCode, w.Code)
		if tc.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.expErr, errMsg)
			continue
		}
		var game model.Game
		readBody(t, w.Body, &game)
		// verify the players are in the game
		for _, p := range g.Players {
			_, ok := game.PlayerColors[p.ID]
			assert.True(t, ok)
		}
	}
}
func TestGinPostCreateInteraction(t *testing.T) {
	testCases := []struct {
		msg  string
		pIDs []string
		req  testRequest
	}{{
		msg: `missing player ID`,
		req: testRequest{
			reqData: network.CreateInteractionRequest{
				PlayerID: ``,
				Mode:     ``,
				Info:     ``,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Needs playerId`,
		},
	}, {
		msg: `bad request body`,
		req: testRequest{
			reqData: struct {
				Field1 string `json:"field1"`
			}{
				Field1: `abc`,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Needs playerId`,
		},
	}, {
		msg: `good request`,
		req: testRequest{
			reqData: network.CreateInteractionRequest{
				PlayerID: `p1`,
				Mode:     `localhost`,
				Info:     ``,
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `unsupported interaction mode`,
		req: testRequest{
			reqData: network.CreateInteractionRequest{
				PlayerID: `p1`,
				Mode:     `abc`,
				Info:     ``,
			},
			expCode: http.StatusBadRequest,
			expErr:  `unsupported interaction mode`,
		},
	}}
	cs, router := newServerAndRouter()
	seedPlayers(t, cs.dbService, 5)
	for _, tc := range testCases {
		// make the request
		body := prepareBody(t, tc.req.reqData)
		w, err := performRequest(router, `POST`, `/create/interaction`, body)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.req.expCode, w.Code)
		_, ok := tc.req.reqData.(network.CreateInteractionRequest)
		if !ok || tc.req.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.req.expErr, errMsg)
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
	testCases := []struct {
		msg    string
		gameID string
		req    testRequest
	}{{
		msg:    `bad game ID`,
		gameID: `123zzz`,
		req: testRequest{
			expCode: http.StatusBadRequest,
			expErr:  `Invalid GameID: strconv.Atoi: parsing "123zzz": invalid syntax`,
		},
	}, {
		msg:    `normal request`,
		gameID: ``,
		req: testRequest{
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg:    `nonexistent game`,
		gameID: `123`,
		req: testRequest{
			expCode: http.StatusNotFound,
			expErr:  `Game not found`,
		},
	}}
	cs, router := newServerAndRouter()
	pIDs := seedPlayers(t, cs.dbService, 2)
	for _, tc := range testCases {
		// seed the db with a game
		g, err := createGame(cs.dbService, pIDs)
		require.NoError(t, err)
		// make the request
		var url string
		if tc.gameID != `` {
			url = `/game/` + tc.gameID
		} else {
			url = `/game/` + fmt.Sprintf(`%d`, g.ID)
		}
		w, err := performRequest(router, `GET`, url, nil)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.req.expCode, w.Code)
		if tc.req.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.req.expErr, errMsg)
			continue
		}
		var game model.Game
		readBody(t, w.Body, &game)
		assert.Equal(t, g.ID, game.ID)
	}
}
func TestGinGetPlayer(t *testing.T) {
	testCases := []struct {
		msg      string
		playerID string
		req      testRequest
	}{{
		msg:      `good request`,
		playerID: `p1`,
		req: testRequest{
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg:      `nonexistent player`,
		playerID: `p9`,
		req: testRequest{
			expCode: http.StatusNotFound,
			expErr:  `Player not found`,
		},
	}}
	cs, router := newServerAndRouter()
	seedPlayers(t, cs.dbService, 2)
	for _, tc := range testCases {
		// make the request
		url := `/player/` + tc.playerID
		w, err := performRequest(router, `GET`, url, nil)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.req.expCode, w.Code)
		if tc.req.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.req.expErr, errMsg)
			continue
		}
		var player model.Player
		readBody(t, w.Body, &player)
		assert.Equal(t, model.PlayerID(tc.playerID), player.ID)
	}
}

// func TestGinPostAction(t *testing.T) {
// 	testCases := []struct {
// 		msg  string
// 		reqs []testRequest
// 	}{{
// 		msg: `good requests`,
// 		reqs: []testRequest{{
// 			reqData: model.PlayerAction{
// 				ID:        `p1`,
// 				Overcomes: model.DealCards,
// 				Action: model.DealAction{
// 					NumShuffles: 1,
// 				},
// 			},
// 			expCode: http.StatusOK,
// 			expErr:  ``,
// 		}},
// 	}}
// 	cs, router := newServerAndRouter()
// 	seedPlayers(t, cs.dbService, 2)
// 	for _, tc := range testCases {
// 		// make the request
// 		url := `/player/` + tc.playerID
// 		w, err := performRequest(router, `GET`, url, nil)
// 		require.NoError(t, err)
// 		// verify
// 		require.Equal(t, tc.req.expCode, w.Code)
// 		if tc.req.expCode != http.StatusOK {
// 			errMsg := readError(t, w)
// 			assert.Equal(t, tc.req.expErr, errMsg)
// 			continue
// 		}
// 		var player model.Player
// 		readBody(t, w.Body, &player)
// 		assert.Equal(t, model.PlayerID(tc.playerID), player.ID)
// 	}
// }
