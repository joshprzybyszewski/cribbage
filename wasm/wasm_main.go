// +build js,wasm

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/callbacks"
)

const (
	serverDomain = `http://localhost:8080`
)

func main() {
	fmt.Println("hello wasm")

	path, err := getCurrentURLPath()
	if err != nil {
		return
	}

	if path == `/` {
		startHomePage()
	} else if isGamePagePath(path) {
		gID, myID, err := parseOutGamePageIDs(path)
		if err != nil {
			return
		}
		startGamePage(gID, myID)
	}

	println(`exiting main`)
}

func getCurrentURLPath() (string, error) {
	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err != nil {
		return ``, err
	}
	return u.Path, nil
}

func startHomePage() {
	println(`startHomePage`)
	var done chan struct{} = make(chan struct{}, 1)

	usernameChangeCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	js.Global().Get("document").
		Call("getElementById", "createUN").
		Call("addEventListener", "change", usernameChangeCB)

	displaynameChangeCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	js.Global().Get("document").
		Call("getElementById", "createDN").
		Call("addEventListener", "change", displaynameChangeCB)

	createUserCB := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("this, args := %+v, %+v\n", this, args)
		return nil
	})

	js.Global().Get("document").
		Call("getElementById", "createForm").
		Call("addEventListener", "submit", createUserCB)

	<-done

	usernameChangeCB.Release()
	displaynameChangeCB.Release()
	createUserCB.Release()
}

var gamePagePathRegex = regexp.MustCompile(`/user/([a-zA-Z0-9_]+)/game/([0-9]+)$`)

func isGamePagePath(path string) bool {
	return gamePagePathRegex.MatchString(path)
}

func parseOutGamePageIDs(path string) (model.GameID, model.PlayerID, error) {
	found := gamePagePathRegex.FindAllStringSubmatch(path, 1)
	chooseFrom := found[0]
	pIDStr := chooseFrom[1]
	gIDStr := chooseFrom[2]
	gIDInt, err := strconv.Atoi(gIDStr)
	if err != nil {
		return model.InvalidGameID, model.InvalidPlayerID, err
	}
	gID := model.GameID(gIDInt)

	return gID, model.PlayerID(pIDStr), nil
}

func startGamePage(gID model.GameID, myID model.PlayerID) {
	println(`startGamePage`)
	var done chan struct{} = make(chan struct{}, 1)

	g, err := requestGame(gID)
	if err != nil {
		return
	}
	callbacks.SetupGamePage(g, myID)
	println(`game setup`)

	<-done
}

func requestGame(gID model.GameID) (model.Game, error) {
	url := fmt.Sprintf("/game/%v", gID)
	respBytes, err := makeRequest(`GET`, url, nil)
	if err != nil {
		return model.Game{}, err
	}

	g, err := jsonutils.UnmarshalGame(respBytes)
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
}

func makeRequest(method, apiURL string, data io.Reader) ([]byte, error) {
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
