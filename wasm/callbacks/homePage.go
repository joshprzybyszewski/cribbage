// +build js,wasm

package callbacks

import (
	"bytes"
	"encoding/json"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/network"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func SetupHomePage() []Releaser {
	var r []Releaser

	rels := addHomePageCallbacks()
	r = append(r, rels...)

	return r
}

func addHomePageCallbacks() []Releaser {
	var r []Releaser

	r = append(r, getListenersForCreateUser()...)

	return r
}

func getListenersForCreateUser() []Releaser {
	var r []Releaser

	doc := dom.GetWindow().Document()

	submitButton := doc.GetElementByID(consts.CreateUserButtonID).(*dom.HTMLButtonElement)
	usernameInput := doc.GetElementByID(consts.CreateUsernameInputID).(*dom.HTMLInputElement)
	displayNameInput := doc.GetElementByID(consts.CreateDisplaynameInputID).(*dom.HTMLInputElement)

	recalcEnabled := func() {
		oldDisabled := submitButton.Disabled()
		newDisabled := len(usernameInput.Value()) == 0 || len(displayNameInput.Value()) == 0
		if newDisabled != oldDisabled {
			submitButton.SetDisabled(newDisabled)
		}
	}

	cb := func(e dom.Event) {
		recalcEnabled()
	}
	r = append(r, getChangeHandlerForID(consts.CreateUsernameInputID, cb))
	r = append(r, getInputHandlerForID(consts.CreateUsernameInputID, cb))
	r = append(r, getChangeHandlerForID(consts.CreateDisplaynameInputID, cb))
	r = append(r, getInputHandlerForID(consts.CreateDisplaynameInputID, cb))

	listener := getClickHandlerForID(consts.CreateUserButtonID, func(e dom.Event) {
		e.PreventDefault()
		username := usernameInput.Value()
		displayname := displayNameInput.Value()
		cpr := network.CreatePlayerRequest{
			Player: network.Player{
				ID:   model.PlayerID(username),
				Name: displayname,
			},
		}

		go func() {
			inputBytes, err := json.Marshal(cpr)
			respBytes, err := actions.MakeRequest(`POST`, `/create/player`, bytes.NewBuffer(inputBytes))
			if err != nil {
				println("Got error on MakeRequest: " + err.Error())
				return
			}
			me := network.CreatePlayerResponse{}
			err = json.Unmarshal(respBytes, &me)
			if err != nil {
				println("Got error on json.Unmarshal CreatePlayerResponse: " + err.Error())
				return
			}
			myUsername := string(me.Player.ID)
			goToPath(`/user/` + myUsername)
		}()
	})

	r = append(r, listener)
	return r
}
