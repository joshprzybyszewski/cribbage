// +build js,wasm

package callbacks

import (
	"honnef.co/go/js/dom/v2"

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
		username := usernameInput.Value()
		displayname := displayNameInput.Value()

		// we might need to wrap this in a go func
		go func() {
			_, err := actions.MakeRequest(`POST`, `/create/player/`+username+`/`+displayname, nil)
			if err != nil {
				println("Got error on MakeRequest: " + err.Error())
				return
			}
		}()
	})

	r = append(r, listener)
	return r
}
