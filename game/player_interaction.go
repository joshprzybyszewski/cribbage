package game

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

type PlayerInteraction interface {
	AskToShuffle() bool
	AskForCribCards(dealerColor model.PlayerColor, desired int, hand []model.Card) []model.Card
	AskForCut() float64
	TellAboutCut(model.Card)
	AskToPeg(hand, prevPegs []model.Card, curPeg int) (c model.Card, sayGo bool)
	TellAboutScores(cur, lag map[model.PlayerColor]int, msgs ...string)
}
