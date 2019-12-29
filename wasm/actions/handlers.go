// +build js,wasm

package actions

import (
	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

// TODO call this when we've clicked a button to submit...
func getCribAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	cards := make([]model.Card, 0, 2)

	// get all of "my cards" that are activated
	elems := dom.GetWindow().Document().QuerySelector(".card.mine." + consts.ActivatedCardClassName)
	for _, elem := range elems {
		htmlElem := elem.(*dom.HTMLDivElement)
		cards = append(cards, model.NewCardFromString(htmlElem.TextContent()))
	}

	bca := model.BuildCribAction{
		Cards: cards,
	}

	return model.PlayerAction{
		GameID:    gID,
		ID:        pID,
		Overcomes: model.CribCard,
		Action:    bca,
	}
}
