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
	fmt.Println(string(bs))
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
	testCases := []struct {
		msg         string
		username    string
		displayName string
		expCode     int
	}{{
		msg:         `normal stuff`,
		username:    `abc`,
		displayName: `def`,
		expCode:     http.StatusOK,
	}, {
		msg:         `username with weird characters shouldn't return 404`,
		username:    `#`,
		displayName: `#`,
		expCode:     http.StatusBadRequest,
	}}
	for _, tc := range testCases {
		// setup a new instance of the server each time to clear the db
		cs := newCribbageServer(memory.New())
		router := cs.NewRouter()

		// make the request
		body := prepareBody(t, network.CreatePlayerRequest{
			Username:    model.PlayerID(tc.username),
			DisplayName: tc.displayName,
		})
		w, err := performRequest(router, `POST`, `/create/player`, body)
		require.NoError(t, err)

		// verify
		assert.Equal(t, tc.expCode, w.Code)
		if tc.expCode == http.StatusOK {
			expPlayer := model.Player{
				ID:   model.PlayerID(tc.username),
				Name: tc.displayName,
			}
			var player model.Player
			readBody(t, w.Body, &player)
			assert.NoError(t, err)
			assert.Equal(t, expPlayer, player)
		}
	}
}
