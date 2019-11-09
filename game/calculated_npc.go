package game

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ PlayerInteraction = (*calcNPCInteraction)(nil)

type calcNPCInteraction struct {
	numShuffles int
	myColor     model.PlayerColor
}

func NewCalcNPC(color model.PlayerColor) Player {
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

func (npc *calcNPCInteraction) AskForCribCards(dealerColor model.PlayerColor, desired int, hand []model.Card) []model.Card {
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

func (npc *calcNPCInteraction) AskToPeg(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (model.Card, bool) {
	return strategy.PegHighestCardNow(hand, prevPegs, curPeg)
}

func (npc *calcNPCInteraction) TellAboutScores(cur, lag map[model.PlayerColor]int, msgs ...string) {}
