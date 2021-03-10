package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

func cribStatsGetter(s model.TossSummary) model.TossStats {
	return s.CribStats
}

// GiveCribHighestPotential gives the crib the highest potential pointed crib
func GiveCribHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, cribStatsGetter, willKeepHighestPotential)
}

// GiveCribLowestPotential gives the crib the lowest potential pointed hand
func GiveCribLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, cribStatsGetter, willKeepLowestPotential)
}
