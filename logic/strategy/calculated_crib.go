package strategy

import (
	"github.com/joshprzybyszewski/cribbage/logic/suggestions"
	"github.com/joshprzybyszewski/cribbage/model"
)

// GiveCribHighestPotential gives the crib the highest potential pointed crib
func GiveCribHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
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
	return getBestPotentialCrib(hand, isBetter)
}

// GiveCribLowestPotential gives the crib the lowest potential pointed hand
func GiveCribLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
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
	return getBestPotentialCrib(hand, isBetter)
}

func getBestPotentialCrib(
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
		if isBetter(prevBest, s.CribStats) {
			bestThrow = bestThrow[:0]
			bestThrow = append(bestThrow, s.Tossed...)
			prevBest = s.CribStats
		}
	}

	return bestThrow, nil
}
