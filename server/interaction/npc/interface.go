package npc

import "github.com/joshprzybyszewski/cribbage/model"

type npcLogic interface {
	addToCrib(model.Game, model.PlayerID, int) []model.Card
	peg([]model.Card, []model.PeggedCard, int) (model.Card, bool)
}
