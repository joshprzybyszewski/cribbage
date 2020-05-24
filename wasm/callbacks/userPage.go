// +build js,wasm

package callbacks

import (
	"bytes"
	"encoding/json"
	"fmt"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func SetupUserPage(myID model.PlayerID) []Releaser {
	var r []Releaser

	rels := addUserPageCallbacks(myID)
	r = append(r, rels...)

	return r
}

func addUserPageCallbacks(myID model.PlayerID) []Releaser {
	var r []Releaser

	r = append(r, getListenersForCreateGame(myID)...)

	return r
}

func getListenersForCreateGame(myID model.PlayerID) []Releaser {
	var r []Releaser

	elem := dom.GetWindow().Document().GetElementByID(consts.CreateGameButtonID)
	createButton := elem.(*dom.HTMLButtonElement)

	elem = dom.GetWindow().Document().GetElementByID(consts.CreateGameOpponentInputID)
	usernameInput := elem.(*dom.HTMLInputElement)

	recalcEnabled := func() {
		oldDisabled := createButton.Disabled()
		newDisabled := len(usernameInput.Value()) == 0
		if newDisabled != oldDisabled {
			createButton.SetDisabled(newDisabled)
		}
	}

	cb := func(e dom.Event) {
		recalcEnabled()
	}
	r = append(r, getChangeHandlerForID(consts.CreateGameOpponentInputID, cb))
	r = append(r, getInputHandlerForID(consts.CreateGameOpponentInputID, cb))

	listener := getClickHandlerForID(consts.CreateGameButtonID, func(e dom.Event) {
		username := usernameInput.Value()
		myUsername := string(myID)
		e.PreventDefault()

		createReq := network.CreateGameRequest{
			PlayerIDs: []model.PlayerID{
				model.PlayerID(username),
				model.PlayerID(myUsername),
			},
		}

		// we might need to wrap this in a go func
		go func() {
			inputBytes, err := json.Marshal(createReq)
			if err != nil {
				println("Got error on json.Marshal: " + err.Error())
				return
			}
			bytes, err := actions.MakeRequest(`POST`, `/create/game`, bytes.NewBuffer(inputBytes))
			if err != nil {
				println("Got error on MakeRequest: " + err.Error())
				return
			}
			g, err := jsonutils.UnmarshalGame(bytes)
			if err != nil {
				println("Got error on UnmarshalGame: " + err.Error())
				return
			}
			gIDStr := fmt.Sprintf("%v", g.ID)
			goToPath(`/user/` + myUsername + `/game/` + gIDStr)
		}()
	})

	r = append(r, listener)
	return r
}
