package game

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/strategy"
)

var _ PlayerInteraction = (*calcNPCInteraction)(nil)

type calcNPCInteraction struct {
	numShuffles int
	myColor     PegColor
}

func NewCalcNPC(color PegColor) Player {
	simple := &calcNPCInteraction{
		myColor: color,
	}
	return newPlayer(simple, `Calculated NPC`, color)
}

func (npc *calcNPCInteraction) AskToShuffle() bool {
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

func (npc *calcNPCInteraction) AskForCribCards(dealerColor PegColor, desired int, hand []model.Card) []model.Card {
	if dealerColor == npc.myColor {
		if rand.Int()%2 == 0 {
			// We might not want this, but it is a form of calculation
			return strategy.KeepHandLowestPotential(desired, hand)
		}
		return strategy.GiveCribHighestPotential(desired, hand)
	}

	if rand.Int()%2 == 0 {
		return strategy.KeepHandHighestPotential(desired, hand)
	}
	return strategy.GiveCribLowestPotential(desired, hand)
}

func (npc *calcNPCInteraction) AskForCut() float64 {
	return rand.Float64()
}

func (npc *calcNPCInteraction) TellAboutCut(model.Card) {}

func (npc *calcNPCInteraction) AskToPeg(hand, prevPegs []model.Card, curPeg int) (model.Card, bool) {
	return strategy.PegHighestCardNow(hand, prevPegs, curPeg)
}

func (npc *calcNPCInteraction) TellAboutScores(cur, lag map[PegColor]int, msgs ...string) {}
