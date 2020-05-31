// +build js,wasm

package callbacks

import (
	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/actions"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func SetupGamePage(g model.Game, myID model.PlayerID) []Releaser {
	var r []Releaser

	listener := getClickHandlerForID(consts.RefreshButtonID, func(e dom.Event) {
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

	enableElemsForPhase(g.Phase)
	newRels := getListenersForPhase(g.ID, myID, g.Phase)
	r = append(r, newRels...)

	return r
}

func enableElemsForPhase(phase model.Phase) {
	ids := []string{}

	switch phase {
	case model.Deal:
		ids = append(ids, consts.DealButtonID)
	case model.BuildCrib:
		ids = append(ids, consts.BuildCribButtonID)
	case model.Cut:
		ids = append(ids, consts.CutInputID)
	case model.Pegging:
		ids = append(ids, consts.PegButtonID)
	case model.Counting:
		ids = append(ids, consts.CountHandPtsInputID)
	case model.CribCounting:
		ids = append(ids, consts.CountCribPtsInputID)
	}

	for _, id := range ids {
		elem := dom.GetWindow().Document().GetElementByID(id)
		if button, ok := elem.(*dom.HTMLButtonElement); ok {
			button.SetDisabled(false)
		} else if input, ok := elem.(*dom.HTMLInputElement); ok {
			input.SetDisabled(false)
		} else {
			println(`could not enable desired id: ` + id)
		}
	}
}

func getListenersForPhase(gID model.GameID, myID model.PlayerID, phase model.Phase) []Releaser {
	var rels []Releaser

	switch phase {
	case model.Deal:
		rels = getDealCallbacks(gID, myID)

	case model.BuildCrib:
		rels = getBuildCribCallbacks(gID, myID)

	case model.Cut:
		rels = getCutCallbacks(gID, myID)

	case model.Pegging:
		rels = getPegCallbacks(gID, myID)

	case model.Counting:
		rels = getCountHandCallbacks(gID, myID)

	case model.CribCounting:
		rels = getCountCribCallbacks(gID, myID)
	}

	return rels
}

func getDealCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add a callback to the deal button
	listener := getClickHandlerForID(consts.DealButtonID, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetDealAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func getBuildCribCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add listeners to my cards to signal activated or not
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine")
	for i := range elems {
		// we need to assign here (not in the for-range)
		// so that each callback has its own memory ref
		elem := elems[i]
		listener := elem.AddEventListener(`click`, false, func(e dom.Event) {
			e.PreventDefault()
			if elem.Class().Contains(consts.ActivatedCardClassName) {
				elem.Class().Remove(consts.ActivatedCardClassName)
			} else {
				elem.Class().Add(consts.ActivatedCardClassName)
			}
		})
		r = append(r, listener)
	}

	// add a callback to the build crib button
	listener := getClickHandlerForID(consts.BuildCribButtonID, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetCribAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func getCutCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add a callback to the cut button
	listener := getEnterKeyHandlerForID(consts.CutInputID, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetCutAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func getPegCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add listeners to my cards to signal activated or not
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine:not(.disabled)")
	for i := range elems {
		// we need to assign here (not in the for-range)
		// so that each callback has its own memory ref
		elem := elems[i]
		listener := elem.AddEventListener(`click`, false, func(e dom.Event) {
			e.PreventDefault()
			if elem.Class().Contains(consts.ActivatedCardClassName) {
				elem.Class().Remove(consts.ActivatedCardClassName)
			} else {
				elem.Class().Add(consts.ActivatedCardClassName)
			}
		})
		r = append(r, listener)
	}

	// add a callback to the cut button
	listener := getClickHandlerForID(consts.PegButtonID, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetPegAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func getCountHandCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add a callback to the cut button
	listener := getEnterKeyHandlerForID(consts.CountHandPtsInputID, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetCountHandAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func getCountCribCallbacks(gID model.GameID, pID model.PlayerID) []Releaser {
	var r []Releaser

	// add a callback to the cut button
	listener := getEnterKeyHandlerForID(consts.CountCribPtsInputID, func(e dom.Event) {
		e.PreventDefault()
		pa := actions.GetCountCribAction(gID, pID)
		sendAction(gID, pa)
	})
	r = append(r, listener)

	return r
}

func sendAction(gID model.GameID, pa model.PlayerAction) {
	go func() {
		err := actions.Send(gID, pa)
		if err != nil {
			println("got error sending action:" + err.Error())
		}
	}()
}
