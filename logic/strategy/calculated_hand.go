package strategy

import (
	"github.com/joshprzybyszewski/cribbage/logic/suggestions"
	"github.com/joshprzybyszewski/cribbage/model"
)

// KeepHandHighestPotential will keep the hand with the highest potential score
func KeepHandHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	isBetter := func(old, new model.TossStats) bool {
		if old == nil {
			return true
		}
		if new == nil {
			return false
		}
		if new.Avg() > old.Avg() {
			return true
		}
		return new.Max() > old.Max()
	}
	return getBestPotentialHand(hand, isBetter)
}

// KeepHandLowestPotential will keep the hand with the lowest potential score
func KeepHandLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	isBetter := func(old, new model.TossStats) bool {
		if old == nil {
			return true
		}
		if new == nil {
			return false
		}
		if new.Avg() < old.Avg() {
			return true
		}
		return new.Min() < old.Min()
	}
	return getBestPotentialHand(hand, isBetter)
}

func getBestPotentialHand(
	hand []model.Card,
	isBetter func(old, new model.TossStats) bool,
) ([]model.Card, error) {
	sums, err := suggestions.GetAllTosses(hand)
	if err != nil {
		return nil, err
	}

	lenDeposit := len(hand) - 4
	bestThrow := make([]model.Card, 0, lenDeposit)
	var prevBest model.TossStats

	for _, s := range sums {
		if isBetter(prevBest, s.HandStats) {
			bestThrow = bestThrow[:0]
			bestThrow = append(bestThrow, s.Tossed...)
			prevBest = s.HandStats
		}
	}

	return bestThrow, nil
}
