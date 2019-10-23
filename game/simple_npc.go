package game

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/strategy"
)

var _ PlayerInteraction = (*simpleNPCInteraction)(nil)

type simpleNPCInteraction struct {
	numShuffles int
	myColor     PegColor
}

func NewSimpleNPC(color PegColor) Player {
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

func (npc *simpleNPCInteraction) AskForCribCards(dealerColor PegColor, desired int, hand []cards.Card) []cards.Card {
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

func (npc *simpleNPCInteraction) TellAboutCut(cards.Card) {}

func (npc *simpleNPCInteraction) AskToPeg(hand, prevPegs []cards.Card, curPeg int) (cards.Card, bool) {
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

func (npc *simpleNPCInteraction) TellAboutScores(cur, lag map[PegColor]int, msgs ...string) {}
