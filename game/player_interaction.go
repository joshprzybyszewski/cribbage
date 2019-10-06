package game

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

type PlayerInteraction interface {
	AskToShuffle() bool
	AskForCribCards(dealerColor PegColor, desired int, hand []cards.Card) []cards.Card
	AskForCut() float64
	TellAboutCut(cards.Card)
	AskToPeg(hand, prevPegs []cards.Card, curPeg int) (c cards.Card, sayGo bool)
	TellAboutPegPoints(int)
}