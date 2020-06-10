// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
	"github.com/joshprzybyszewski/cribbage/wasm/callbacks"
)

func main() {
	println(`starting gowasm`)
	defer func() {
		println(`exiting gowasm`)
	}()

	path := getCurrentURLPath()

	var done chan struct{} = make(chan struct{}, 1)
	var rels []callbacks.Releaser

	if path == `/` {
		rels = callbacks.SetupHomePage()
	} else if isUserPagePath(path) {
		myID := parseOutUserPageIDs(path)
		rels = callbacks.SetupUserPage(myID)
	} else if isGamePagePath(path) {
		gID, myID, err := parseOutGamePageIDs(path)
		if err != nil {
			println(`error on parseOutGamePageIDs: ` + err.Error())
			return
		}
		rels, err = startGamePage(gID, myID)
		if err != nil {
			println(`error on startGamePage: ` + err.Error())
			return
		}
	} else {
		println(`unsupported page`)
		return
	}

	<-done

	for _, r := range rels {
		r.Release()
	}

	println(`exiting main`)
}

func getCurrentURLPath() string {
	return strings.TrimPrefix(dom.GetWindow().Location().Pathname(), `/wasm`)
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

var userPagePathRegex = regexp.MustCompile(`/user/([a-zA-Z0-9_]+)$`)

func isUserPagePath(path string) bool {
	return userPagePathRegex.MatchString(path)
}

func parseOutUserPageIDs(path string) model.PlayerID {
	found := userPagePathRegex.FindAllStringSubmatch(path, 1)
	chooseFrom := found[0]
	pIDStr := chooseFrom[1]

	return model.PlayerID(pIDStr)
}

func startGamePage(gID model.GameID, myID model.PlayerID) ([]callbacks.Releaser, error) {
	println(`startGamePage`)

	g, err := requestGame(gID, myID)
	if err != nil {
		println("got error requesting game: " + err.Error())
		return nil, err
	}

	rels := callbacks.SetupGamePage(g, myID)

	return rels, nil
}

func requestGame(gID model.GameID, myID model.PlayerID) (model.Game, error) {
	url := fmt.Sprintf("/game/%v?player=%s", gID, myID)
	respBytes, err := actions.MakeRequest(`GET`, url, nil)
	if err != nil {
		return model.Game{}, err
	}

	ggr := network.GetGameResponse{}
	err = json.Unmarshal(respBytes, &ggr)
	if err != nil {
		return model.Game{}, err
	}

	mg := network.ConvertFromGetGameResponse(ggr)

	return mg, nil
}
