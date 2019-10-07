package game

import (
	"fmt"
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var _ PlayerInteraction = (*dumbNPCInteraction)(nil)

type dumbNPCInteraction struct {
}

func NewDumbNPC(color PegColor) Player {
	dumb := &dumbNPCInteraction{}
	return newPlayer(dumb, `dumb NPC`, color)
}

func (npc *dumbNPCInteraction) AskToShuffle() bool {
	return rand.Intn(100) < 50
}

func (npc *dumbNPCInteraction) AskForCribCards(dealerColor PegColor, desired int, hand []cards.Card) []cards.Card {
	c := hand[0:2]

	return c
}

func (npc *dumbNPCInteraction) AskForCut() float64 {
	return rand.Float64()
}

func (npc *dumbNPCInteraction) TellAboutCut(cards.Card) {}

func (npc *dumbNPCInteraction) AskToPeg(hand, prevPegs []cards.Card, curPeg int) (cards.Card, bool) {
	maxVal := maxPeggingValue - curPeg
	for _, c := range hand {
		if c.PegValue() > maxVal {
			continue
		}
		fmt.Printf("Dumb NPC pegged %+v\n", c.String())
		return c, false
	}

	fmt.Printf("Dumb NPC says go\n")
	return cards.Card{}, true
}

func (npc *dumbNPCInteraction) TellAboutPegPoints(n int) {
	fmt.Printf("Dump NPC received %d points for pegging\n", n)
}