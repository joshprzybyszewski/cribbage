// +build js,wasm

package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
)

func Send(gID model.GameID, pa model.PlayerAction) error {
	b, err := json.Marshal(pa)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)

	url := fmt.Sprintf("/action")

	_, err = MakeRequest(`POST`, url, buf)
	return err
}

func getServerDomain() string {
	loc := dom.GetWindow().Location()
	return loc.Protocol() + `//` + loc.Host()
}

func MakeRequest(method, apiURL string, data io.Reader) ([]byte, error) {
	urlStr := getServerDomain() + apiURL
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
			return nil, errors.New(`bad response: "` + string(bytes) + `"`)
		}

		return nil, errors.New(`bad response from server`)
	} else if err != nil {
		return nil, err
	}

	return bytes, nil
}
