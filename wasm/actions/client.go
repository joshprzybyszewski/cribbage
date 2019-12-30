// +build js,wasm

package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/joshprzybyszewski/cribbage/model"
)

const (
	serverDomain = `http://localhost:8080`
)

func Send(gID model.GameID, pa model.PlayerAction) error {
	b, err := json.Marshal(pa)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)

	url := fmt.Sprintf("/action/%d", gID)

	fmt.Printf("url: %+v\n", url)

	_, err = MakeRequest(`POST`, url, buf)
	return err
}

func MakeRequest(method, apiURL string, data io.Reader) ([]byte, error) {
	urlStr := serverDomain + apiURL
	req, err := http.NewRequest(method, urlStr, data)
	if err != nil {
		return nil, err
	}

	server := &http.Client{}
	response, err := server.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		contentType := response.Header.Get("Content-Type")
		if strings.Contains(contentType, `text/plain`) {
			return nil, fmt.Errorf("bad response: \"%s\"", string(bytes))
		}

		return nil, fmt.Errorf("bad response from server")
	} else if err != nil {
		return nil, err
	}

	return bytes, nil
}
