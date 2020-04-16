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

func TestGinPostCreatePlayer(t *testing.T) {
	testCases := []struct {
		msg     string
		data    CreatePlayerData
		expRes  model.Player
		expCode int
	}{
		/*
			TODO add test cases once the server is more testable. We need to rewrite it to inject persistence in so we can mock it here
			(or at least set it to memory instead of mongo by default)
		*/
	}
	cs := &cribbageServer{}
	router := cs.NewRouter()
	for _, tc := range testCases {
		reqBytes, err := json.Marshal(tc.data)
		require.NoError(t, err)
		body := bytes.NewReader(reqBytes)
		w, err := performRequest(router, `POST`, `/create/player`, body)
		require.NoError(t, err)
		assert.Equal(t, tc.expCode, w.Code)
		if tc.expCode == http.StatusOK {
			var player model.Player
			readBody(t, w.Body, &player)
			assert.NoError(t, err)
			assert.Equal(t, tc.expRes, player)
		}
	}
}
