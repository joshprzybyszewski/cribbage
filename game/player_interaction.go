package game

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type PlayerInteraction interface {
	AskToShuffle() bool
	AskForCribCards(dealerColor PegColor, desired int, hand []model.Card) []model.Card
	AskForCut() float64
	TellAboutCut(model.Card)
	AskToPeg(hand, prevPegs []model.Card, curPeg int) (c model.Card, sayGo bool)
	TellAboutScores(cur, lag map[PegColor]int, msgs ...string)
}
