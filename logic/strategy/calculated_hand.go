package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

// KeepHandHighestPotential will keep the hand with the highest potential score
func KeepHandHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, newTossEvaluator(true, highestIsBetter))
}

// KeepHandLowestPotential will keep the hand with the lowest potential score
func KeepHandLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, newTossEvaluator(true, lowestIsBetter))
}
