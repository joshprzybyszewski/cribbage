package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ handEvaluator = (*highestCribEvaluator)(nil)

type highestCribEvaluator struct{}

func (*highestCribEvaluator) getStats(ts model.TossSummary) model.TossStats {
	return ts.CribStats
}

func (*highestCribEvaluator) isBetter(old, new model.TossStats) bool {
	return willKeepHighestPotential(old, new)
}

// GiveCribHighestPotential gives the crib the highest potential pointed crib
func GiveCribHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, &highestCribEvaluator{})
}

var _ handEvaluator = (*lowestCribEvaluator)(nil)

type lowestCribEvaluator struct{}

func (*lowestCribEvaluator) getStats(ts model.TossSummary) model.TossStats {
	return ts.CribStats
}

func (*lowestCribEvaluator) isBetter(old, new model.TossStats) bool {
	return willKeepLowestPotential(old, new)
}

// GiveCribLowestPotential gives the crib the lowest potential pointed hand
func GiveCribLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, &lowestCribEvaluator{})
}
