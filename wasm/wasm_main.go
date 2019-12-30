// +build js,wasm

package main

import (
	"fmt"
	"regexp"
	"strconv"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
	"github.com/joshprzybyszewski/cribbage/wasm/callbacks"
)

func main() {
	fmt.Println("hello wasm")

	path := getCurrentURLPath()

	var done chan struct{} = make(chan struct{}, 1)
	var rels []callbacks.Releaser

	if path == `/` {
		rels = callbacks.SetupHomePage()
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

	elem := dom.GetWindow().Document().GetElementByID("kill")
	if htmlElem, ok := elem.(*dom.HTMLDivElement); ok {
		rels = append(rels, htmlElem.AddEventListener(`click`, false, func(e dom.Event) {
			done <- struct{}{}
		}))
	}

	<-done

	for _, r := range rels {
		r.Release()
	}

	println(`exiting main`)
}

func getCurrentURLPath() string {
	return dom.GetWindow().Location().Pathname()
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

func startGamePage(gID model.GameID, myID model.PlayerID) ([]callbacks.Releaser, error) {
	g, err := requestGame(gID)
	if err != nil {
		return nil, err
	}
	rels := callbacks.SetupGamePage(g, myID)
	return rels, nil
}

func requestGame(gID model.GameID) (model.Game, error) {
	url := fmt.Sprintf("/game/%v", gID)
	respBytes, err := actions.MakeRequest(`GET`, url, nil)
	if err != nil {
		return model.Game{}, err
	}

	g, err := jsonutils.UnmarshalGame(respBytes)
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
}
