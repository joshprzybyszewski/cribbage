package npc

import "github.com/joshprzybyszewski/cribbage/model"

type npcLogic interface {
	addToCrib(hand []model.Card, isDealer bool) []model.Card
	peg([]model.Card, []model.PeggedCard, int) (model.Card, bool)
}
