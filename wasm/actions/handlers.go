// +build js,wasm

package actions

import (
	"honnef.co/go/js/dom/v2"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
	"github.com/joshprzybyszewski/cribbage/wasm/consts"
)

func GetDealAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	da := model.DealAction{
		NumShuffles: rand.Intn(10),
	}

	return model.PlayerAction{
		GameID:    gID,
		ID:        pID,
		Overcomes: model.DealCards,
		Action:    da,
	}
}

func GetCribAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	cards := make([]model.Card, 0, 2)

	// get all of "my cards" that are activated
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine." + consts.ActivatedCardClassName)
	for _, elem := range elems {
		cards = append(cards, model.NewCardFromString(elem.ID()))
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

func GetCutAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	elem := dom.GetWindow().Document().GetElementByID(consts.CutInputID)
	input := elem.(*dom.HTMLInputElement)
	perc := input.ValueAsNumber()
	if perc >= 1.0 && perc <= 100.0 {
		// assume it's entered in hundredths
		perc = perc / 100.0
	}

	cda := model.CutDeckAction{
		Percentage: perc,
	}

	return model.PlayerAction{
		GameID:    gID,
		ID:        pID,
		Overcomes: model.CutCard,
		Action:    cda,
	}
}

func GetPegAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	cards := make([]model.Card, 0, 1)

	// get all of "my cards" that are activated
	elems := dom.GetWindow().Document().QuerySelectorAll(".card.mine." + consts.ActivatedCardClassName)
	for _, elem := range elems {
		cards = append(cards, model.NewCardFromString(elem.ID()))
	}

	pegA := model.PegAction{}

	if len(cards) == 1 {
		pegA.Card = cards[0]
	} else {
		pegA.SayGo = true
	}

	return model.PlayerAction{
		GameID:    gID,
		ID:        pID,
		Overcomes: model.PegCard,
		Action:    pegA,
	}
}

func GetCountHandAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	elem := dom.GetWindow().Document().GetElementByID(consts.CountHandPtsInputID)
	input := elem.(*dom.HTMLInputElement)
	pts := int(input.ValueAsNumber())

	cha := model.CountHandAction{
		Pts: pts,
	}

	return model.PlayerAction{
		GameID:    gID,
		ID:        pID,
		Overcomes: model.CountHand,
		Action:    cha,
	}
}

func GetCountCribAction(gID model.GameID, pID model.PlayerID) model.PlayerAction {
	elem := dom.GetWindow().Document().GetElementByID(consts.CountCribPtsInputID)
	input := elem.(*dom.HTMLInputElement)
	pts := int(input.ValueAsNumber())

	cca := model.CountCribAction{
		Pts: pts,
	}

	return model.PlayerAction{
		GameID:    gID,
		ID:        pID,
		Overcomes: model.CountCrib,
		Action:    cca,
	}
}
