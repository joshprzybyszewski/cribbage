// +build js,wasm

package callbacks

import (
	"encoding/json"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
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

	elem := dom.GetWindow().Document().GetElementByID(consts.CreateUserButtonID)
	submitButton := elem.(*dom.HTMLButtonElement)

	elem = dom.GetWindow().Document().GetElementByID(consts.CreateUsernameInputID)
	usernameInput := elem.(*dom.HTMLInputElement)
	elem = dom.GetWindow().Document().GetElementByID(consts.CreateDisplaynameInputID)
	displayNameInput := elem.(*dom.HTMLInputElement)

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
		// TODO for some reason the username and displayname are the same
		username := usernameInput.Value()
		displayname := displayNameInput.Value()

		go func() {
			bytes, err := actions.MakeRequest(`POST`, `/create/player/`+username+`/`+displayname, nil)
			if err != nil {
				println("Got error on MakeRequest: " + err.Error())
				return
			}
			me := model.Player{}
			json.Unmarshal(bytes, &me)
			myUsername := string(me.ID)
			goToPath(`/user/` + myUsername)
		}()
	})

	r = append(r, listener)
	return r
}
