package interaction

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/joshprzybyszewski/cribbage/model"
)

var _ Player = (*localhostPlayer)(nil)

type localhostPlayer struct {
	pID  model.PlayerID
	port int
}

func newLocalhostPlayer(pID model.PlayerID, info interface{}) *localhostPlayer {
	pStr, ok := info.(string)
	if !ok {
		pStr = `8081`
	}
	port, err := strconv.Atoi(pStr)
	if err != nil {
		port = 8082
	}
	return &localhostPlayer{
		pID:  pID,
		port: port,
	}
}

func (lhp *localhostPlayer) ID() model.PlayerID {
	return lhp.pID
}

func (lhp *localhostPlayer) NotifyBlocking(b model.Blocker, g model.Game, i string) error {
	return lhp.notify(`blocking`, nil)
}

func (lhp *localhostPlayer) NotifyMessage(g model.Game, msg string) error {
	return lhp.notify(`message`, nil)
}

func (lhp *localhostPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return lhp.notify(`score`, nil)
}

func (lhp *localhostPlayer) notify(endpoint string, data io.ReadCloser) error {
	urlStr := fmt.Sprintf("localhost:%d/%s", lhp.port, endpoint)
	url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	client := &http.Client{}

	request := http.Request{
		Method: `POST`,
		URL:    url,
		Body:   data,
	}
	response, err := client.Do(&request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(`bad response: %+v`, response)
	}

	_, err = ioutil.ReadAll(response.Body)
	return err
}
