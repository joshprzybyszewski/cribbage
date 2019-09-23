package game

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var _ Player = (*dumbNPC)(nil)

type dumbNPC struct {
	*player
}

func NewDumbNPC(color PegColor) Player {
	return &dumbNPC{
		player: newPlayer(`dumb NPC`, color),
	}
}

func (npc *dumbNPC) Shuffle() {
	if !npc.IsDealer() {
		return
	}

	npc.deck.Shuffle()
	for rand.Intn(100) < 50 {
		npc.deck.Shuffle()
	}
}

func (npc *dumbNPC) AddToCrib() []cards.Card {
	c := npc.hand[0:2]

	npc.hand = npc.hand[2:]

	return c
}

func (npc *dumbNPC) Cut() float64 {
	return rand.Float64()
}

func (npc *dumbNPC) Peg(maxVal int) (cards.Card, bool, bool) {
	if len(npc.pegged) == 4 {
		return cards.Card{}, false, false
	}

	cardToPlay := cards.Card{}
	sayGo := true

	for _, c := range npc.hand {
		if _, ok := npc.pegged[c]; ok {
			continue
		}
		if c.PegValue() > maxVal {
			continue
		}
		cardToPlay = c
		sayGo = false
		break
	}

	if sayGo {
		return cards.Card{}, true, true
	}

	return cardToPlay, false, true
}
