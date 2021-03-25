package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

// GiveCribHighestPotential gives the crib the highest potential pointed crib
func GiveCribHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, newTossEvaluator(false, highestIsBetter))
}

// GiveCribLowestPotential gives the crib the lowest potential pointed hand
func GiveCribLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, newTossEvaluator(false, lowestIsBetter))
}
