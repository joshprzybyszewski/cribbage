// +build js,wasm

package callbacks

import (
	"fmt"

	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func SetupGamePage(g model.Game, myID model.PlayerID) []Releaser {
	var r []Releaser

	elem := dom.GetWindow().Document().GetElementByID("refreshGame")
	button := elem.(*dom.HTMLButtonElement)
	listener := button.AddEventListener(`click`, false, func(e dom.Event) {
		e.PreventDefault()
		// TODO this is terribly inefficient. We should have websockets on the game page
		// so that we can get updates for just what we need. But since this is POC, let's
		// do something hacky
		dom.GetWindow().Location().Call(`reload`, true)
	})
	r = append(r, listener)

	if _, ok := g.BlockingPlayers[myID]; !ok {
		// Nothing to do when I'm not blocking
		// TODO this may not always be true
		return r
	}

	switch g.Phase {
	case model.BuildCrib:
		enableBuildCribElems()
		rels := addBuildCribCallbacks(g.ID, myID)
		r = append(r, rels...)
	}

	return r
}

func enableBuildCribElems() {
	elem := dom.GetWindow().Document().GetElementByID("buildCribButton")
	button := elem.(*dom.HTMLButtonElement)
	button.SetDisabled(false)
}

func addBuildCribCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add listeners to my cards to signal activated or not
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine")
	for _, elem := range elems {
		htmlElem := elem.(*dom.HTMLDivElement)
		listener := htmlElem.AddEventListener(`click`, false, func(e dom.Event) {
			e.PreventDefault()
			if htmlElem.Class().Contains(consts.ActivatedCardClassName) {
				htmlElem.Class().Remove(consts.ActivatedCardClassName)
			} else {
				htmlElem.Class().Add(consts.ActivatedCardClassName)
			}
		})
		r = append(r, listener)
	}

	// add a callback to the build crib button
	elem := dom.GetWindow().Document().GetElementByID("buildCribButton")
	button := elem.(*dom.HTMLButtonElement)
	listener := button.AddEventListener(`click`, false, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetCribAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func sendAction(gID model.GameID, pa model.PlayerAction) {
	go func() {
		err := actions.Send(gID, pa)
		if err != nil {
			fmt.Printf("got error sending action: %+v\n", err)
		}
	}()
}
