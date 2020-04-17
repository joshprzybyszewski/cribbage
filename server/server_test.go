package server

import (
	"bytes"
	"encoding/json"
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

func prepareBody(t *testing.T, v interface{}) io.Reader {
	reqBytes, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(reqBytes)
}

func TestGinPostCreatePlayer(t *testing.T) {
	type testRequest struct {
		username    string
		displayName string
		expCode     int
		expErr      string
	}

	testCases := []struct {
		msg  string
		reqs []testRequest
	}{{
		msg: `normal stuff`,
		reqs: []testRequest{{
			username:    `abc`,
			displayName: `def`,
			expCode:     http.StatusOK,
			expErr:      ``,
		}},
	}, {
		msg: `username with weird characters shouldn't return 404`,
		reqs: []testRequest{{
			username:    `#`,
			displayName: `#`,
			expCode:     http.StatusBadRequest,
			expErr:      `Username must be alphanumeric`,
		}},
	}, {
		msg: `creating the same player errors`,
		reqs: []testRequest{{
			username:    `abc`,
			displayName: `def`,
			expCode:     http.StatusOK,
			expErr:      ``,
		}, {
			username:    `abc`,
			displayName: `def`,
			expCode:     http.StatusBadRequest,
			expErr:      `Username already exists`,
		}},
	}}
	for _, tc := range testCases {
		// setup a new instance of the server each time to clear the db
		cs := newCribbageServer(memory.New())
		router := cs.NewRouter()

		// make the requests
		for _, r := range tc.reqs {
			body := prepareBody(t, network.CreatePlayerRequest{
				Username:    model.PlayerID(r.username),
				DisplayName: r.displayName,
			})
			w, err := performRequest(router, `POST`, `/create/player`, body)
			require.NoError(t, err)
			// verify
			require.Equal(t, r.expCode, w.Code)
			if r.expCode == http.StatusOK {
				expPlayer := model.Player{
					ID:   model.PlayerID(r.username),
					Name: r.displayName,
				}
				var player model.Player
				readBody(t, w.Body, &player)
				assert.NoError(t, err)
				assert.Equal(t, expPlayer, player)
			} else {
				errMsgBytes, err := ioutil.ReadAll(w.Body)
				require.NoError(t, err)
				assert.Equal(t, r.expErr, string(errMsgBytes))
			}
		}
		memory.Clear()
	}
}
