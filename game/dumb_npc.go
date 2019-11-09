package game

import (
	"fmt"
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/model"
)

var _ PlayerInteraction = (*dumbNPCInteraction)(nil)

type dumbNPCInteraction struct {
	numShuffles int
}

func NewDumbNPC(color PegColor) Player {
	dumb := &dumbNPCInteraction{}
	return newPlayer(dumb, `dumb NPC`, color)
}

func (npc *dumbNPCInteraction) AskToShuffle() bool {
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

func (npc *dumbNPCInteraction) AskForCribCards(dealerColor PegColor, desired int, hand []model.Card) []model.Card {
	c := hand[0:2]

	return c
}

func (npc *dumbNPCInteraction) AskForCut() float64 {
	return rand.Float64()
}

func (npc *dumbNPCInteraction) TellAboutCut(model.Card) {}

func (npc *dumbNPCInteraction) AskToPeg(hand, prevPegs []model.Card, curPeg int) (model.Card, bool) {
	maxVal := maxPeggingValue - curPeg
	for _, c := range hand {
		if c.PegValue() > maxVal {
			continue
		}
		fmt.Printf("Dumb NPC pegged %+v\n", c.String())
		return c, false
	}

	fmt.Printf("Dumb NPC says go\n")
	return model.Card{}, true
}

func (npc *dumbNPCInteraction) TellAboutScores(cur, lag map[PegColor]int, msgs ...string) {}
