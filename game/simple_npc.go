package game

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
)

var _ PlayerInteraction = (*simpleNPCInteraction)(nil)

type simpleNPCInteraction struct {
	numShuffles int
	myColor     model.PlayerColor
}

func NewSimpleNPC(color model.PlayerColor) Player {
	simple := &simpleNPCInteraction{
		myColor: color,
	}
	return newPlayer(simple, `Simple NPC`, color)
}

func (npc *simpleNPCInteraction) AskToShuffle() bool {
	npc.numShuffles++

	if npc.numShuffles <= 1 {
		return true
	}

	shouldContinue := rand.Intn(100) < npc.numShuffles
	if !shouldContinue {
		npc.numShuffles = 0
	}

	return shouldContinue
}

func (npc *simpleNPCInteraction) AskForCribCards(dealerColor model.PlayerColor, desired int, hand []model.Card) []model.Card {
	if dealerColor == npc.myColor {
		if rand.Int()%2 == 0 {
			return strategy.GiveCribFifteens(desired, hand)
		}
		return strategy.GiveCribPairs(desired, hand)
	}

	if rand.Int()%2 == 0 {
		return strategy.AvoidCribFifteens(desired, hand)
	}
	return strategy.AvoidCribPairs(desired, hand)
}

func (npc *simpleNPCInteraction) AskForCut() float64 {
	return rand.Float64()
}

func (npc *simpleNPCInteraction) TellAboutCut(model.Card) {}

func (npc *simpleNPCInteraction) AskToPeg(hand, prevPegs []model.Card, curPeg int) (model.Card, bool) {
	switch rand.Int() % 4 {
	case 0:
		return strategy.PegToFifteen(hand, prevPegs, curPeg)
	case 1:
		return strategy.PegToThirtyOne(hand, prevPegs, curPeg)
	case 2:
		return strategy.PegToPair(hand, prevPegs, curPeg)
	}
	return strategy.PegToRun(hand, prevPegs, curPeg)
}

func (npc *simpleNPCInteraction) TellAboutScores(cur, lag map[model.PlayerColor]int, msgs ...string) {}
