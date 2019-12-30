// +build js,wasm

package callbacks

import (
	"encoding/json"
	"fmt"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
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

	elem := dom.GetWindow().Document().GetElementByID("create-user-submit")
	submitButton := elem.(*dom.HTMLInputElement)

	elem = dom.GetWindow().Document().GetElementByID("createUN")
	usernameInput := elem.(*dom.HTMLInputElement)
	elem = dom.GetWindow().Document().GetElementByID("createDN")
	displayNameInput := elem.(*dom.HTMLInputElement)

	recalcEnabled := func() {
		oldDisabled := submitButton.Disabled()
		newDisabled := usernameInput.Value() == `` || displayNameInput.Value() == ``
		if newDisabled != oldDisabled {
			submitButton.SetDisabled(newDisabled)
		}
	}

	listener := usernameInput.AddEventListener(`change`, false, func(e dom.Event) {
		recalcEnabled()
	})
	r = append(r, listener)
	listener = usernameInput.AddEventListener(`input`, false, func(e dom.Event) {
		recalcEnabled()
	})
	r = append(r, listener)
	listener = displayNameInput.AddEventListener(`change`, false, func(e dom.Event) {
		recalcEnabled()
	})
	r = append(r, listener)
	listener = displayNameInput.AddEventListener(`input`, false, func(e dom.Event) {
		recalcEnabled()
	})
	r = append(r, listener)

	elem = dom.GetWindow().Document().GetElementByID("createForm")
	form := elem.(*dom.HTMLFormElement)
	listener = form.AddEventListener(`submit`, false, func(e dom.Event) {
		username := usernameInput.Value()
		displayname := displayNameInput.Value()

		respBytes, err := actions.MakeRequest(`POST`, `/create/player/`+username+`/`+displayname, nil)
		if err != nil {
			fmt.Printf("Got error on MakeRequest: %v\n", err)
			return
		}

		player := model.Player{}
		err = json.Unmarshal(respBytes, &player)
		if err != nil {
			fmt.Printf("Got error on Unmarshal: %v\n", err)
			return
		}

		fmt.Printf("Your player ID is: %v\n", player.ID)
	})
	r = append(r, listener)
	return r
}
