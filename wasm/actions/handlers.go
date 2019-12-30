// +build js,wasm

package actions

import (
	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func GetCribAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	cards := make([]model.Card, 0, 2)

	// get all of "my cards" that are activated
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine." + consts.ActivatedCardClassName)
	for _, elem := range elems {
		div := elem.(*dom.HTMLDivElement)
		cards = append(cards, model.NewCardFromString(div.ID()))
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
