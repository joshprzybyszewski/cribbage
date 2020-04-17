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

func TestGinPostCreatePlayer(t *testing.T) {
	type testRequest struct {
		reqData interface{}
		expCode int
		expErr  string
	}

	testCases := []struct {
		msg  string
		reqs []testRequest
	}{{
		msg: `normal stuff`,
		reqs: []testRequest{{
			reqData: network.CreatePlayerRequest{
				Username:    `abc`,
				DisplayName: `def`,
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}},
	}, {
		msg: `username with weird characters shouldn't return 404`,
		reqs: []testRequest{{
			reqData: network.CreatePlayerRequest{
				Username:    `#`,
				DisplayName: `#`,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username must be alphanumeric`,
		}},
	}, {
		msg: `creating the same player errors`,
		reqs: []testRequest{{
			reqData: network.CreatePlayerRequest{
				Username:    `abc`,
				DisplayName: `def`,
			},
			expCode: http.StatusOK,
			expErr:  ``,
		}, {
			reqData: network.CreatePlayerRequest{
				Username:    `abc`,
				DisplayName: `def`,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username already exists`,
		}},
	}, {
		msg: `empty username`,
		reqs: []testRequest{{
			reqData: network.CreatePlayerRequest{
				Username:    ``,
				DisplayName: `def`,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username is required`,
		}},
	}, {
		msg: `empty display name`,
		reqs: []testRequest{{
			reqData: network.CreatePlayerRequest{
				Username:    `abc`,
				DisplayName: ``,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Display name is required`,
		}},
	}, {
		msg: `send wrong JSON data`,
		reqs: []testRequest{{
			reqData: struct {
				Field1 string `json:"field1"`
			}{
				Field1: `abc`,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Username is required`,
		}},
	}}
	for _, tc := range testCases {
		_, router := newServerAndRouter()

		// make the requests
		for _, r := range tc.reqs {
			body := prepareBody(t, r.reqData)
			w, err := performRequest(router, `POST`, `/create/player`, body)
			require.NoError(t, err)
			// verify
			require.Equal(t, r.expCode, w.Code)
			cpr, ok := r.reqData.(network.CreatePlayerRequest)
			if !ok || r.expCode != http.StatusOK {
				errMsg := readError(t, w)
				assert.Equal(t, r.expErr, errMsg)
				continue
			}
			expPlayer := model.Player{
				ID:   model.PlayerID(cpr.Username),
				Name: cpr.DisplayName,
			}
			var player model.Player
			readBody(t, w.Body, &player)
			assert.NoError(t, err)
			assert.Equal(t, expPlayer, player)
		}
	}
}
func TestGinPostCreateGame(t *testing.T) {
	type testRequest struct {
		reqData interface{}
		expCode int
		expErr  string
	}

	testCases := []struct {
		msg  string
		pIDs []string
		req  testRequest
	}{{
		msg: `two player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `three player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `p3`},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `four player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `p3`, `p4`},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `one player game is an error`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 1`,
		},
	}, {
		msg: `five player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `p3`, `p4`, `p5`},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 5`,
		},
	}, {
		msg: `zero player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 0`,
		},
	}, {
		msg: `missing player id`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, ``, `p4`},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid player ID at index 2`,
		},
	}, {
		msg: `invalid player id`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `#`, `p4`},
			},
			expCode: http.StatusInternalServerError,
			expErr:  `createGame error: player not found`,
		},
	}, {
		msg: `create a game with nonexistent players`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p6`},
			},
			expCode: http.StatusInternalServerError,
			expErr:  `createGame error: player not found`,
		},
	}, {
		msg: `bad request body`,
		req: testRequest{
			reqData: struct {
				Field1 string `json:"Field1"`
			}{
				Field1: `abc`,
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 0`,
		},
	}}
	cs, router := newServerAndRouter()
	// seed the db with players
	for i := 0; i < 5; i++ {
		_, err := createPlayer(cs.dbService, fmt.Sprintf(`p%d`, i+1), `name`)
		require.NoError(t, err)
	}
	for _, tc := range testCases {

		// make the request
		body := prepareBody(t, tc.req.reqData)
		w, err := performRequest(router, `POST`, `/create/game`, body)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.req.expCode, w.Code)
		cgr, ok := tc.req.reqData.(network.CreateGameRequest)
		if !ok || tc.req.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.req.expErr, errMsg)
			continue
		}
		var game model.Game
		readBody(t, w.Body, &game)
		// verify the players are in the game
		for _, pID := range cgr.PlayerIDs {
			_, ok := game.PlayerColors[model.PlayerID(pID)]
			assert.True(t, ok)
		}
	}
}
func TestGinPostCreateInteraction(t *testing.T) {
	type testRequest struct {
		reqData interface{}
		expCode int
		expErr  string
	}

	testCases := []struct {
		msg  string
		pIDs []string
		req  testRequest
	}{{
		msg: `two player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `three player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `p3`},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `four player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `p3`, `p4`},
			},
			expCode: http.StatusOK,
			expErr:  ``,
		},
	}, {
		msg: `one player game is an error`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 1`,
		},
	}, {
		msg: `five player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `p3`, `p4`, `p5`},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 5`,
		},
	}, {
		msg: `zero player game`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid num players: 0`,
		},
	}, {
		msg: `missing player id`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, ``, `p4`},
			},
			expCode: http.StatusBadRequest,
			expErr:  `Invalid player ID at index 2`,
		},
	}, {
		msg: `invalid player id`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p2`, `#`, `p4`},
			},
			expCode: http.StatusInternalServerError,
			expErr:  `createGame error: player not found`,
		},
	}, {
		msg: `create a game with nonexistent players`,
		req: testRequest{
			reqData: network.CreateGameRequest{
				PlayerIDs: []string{`p1`, `p6`},
			},
			expCode: http.StatusInternalServerError,
			expErr:  `createGame error: player not found`,
		},
	}}
	cs, router := newServerAndRouter()
	// seed the db with players
	for i := 0; i < 5; i++ {
		_, err := createPlayer(cs.dbService, fmt.Sprintf(`p%d`, i+1), `name`)
		require.NoError(t, err)
	}
	for _, tc := range testCases {

		// make the request
		body := prepareBody(t, tc.req.reqData)
		w, err := performRequest(router, `POST`, `/create/game`, body)
		require.NoError(t, err)
		// verify
		require.Equal(t, tc.req.expCode, w.Code)
		cgr, ok := tc.req.reqData.(network.CreateGameRequest)
		if !ok || tc.req.expCode != http.StatusOK {
			errMsg := readError(t, w)
			assert.Equal(t, tc.req.expErr, errMsg)
			continue
		}
		var game model.Game
		readBody(t, w.Body, &game)
		// verify the players are in the game
		for _, pID := range cgr.PlayerIDs {
			_, ok := game.PlayerColors[model.PlayerID(pID)]
			assert.True(t, ok)
		}
	}
}
