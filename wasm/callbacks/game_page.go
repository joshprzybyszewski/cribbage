// +build js,wasm

package callbacks

import (
	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func SetupGamePage(g model.Game, myID model.PlayerID) {
	if _, ok := g.BlockingPlayers[myID]; !ok {
		// Nothing to do when I'm not blocking
		return
	}

	switch g.Phase {
	case model.BuildCrib:
		addBuildCribCallbacks()
	}
}

func addBuildCribCallbacks() {
	// add listeners to my cards to signal activated or not
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine")
	for _, elem := range elems {
		htmlElem := elem.(*dom.HTMLDivElement)
		htmlElem.AddEventListener(`click`, false, func(e dom.Event) {
			e.PreventDefault()
			if htmlElem.Class().Contains(consts.ActivatedCardClassName) {
				htmlElem.Class().Remove(consts.ActivatedCardClassName)
			} else {
				htmlElem.Class().Add(consts.ActivatedCardClassName)
			}
		})
	}
}
