package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ handEvaluator = (*highestHandEvaluator)(nil)

type highestHandEvaluator struct{}

func (*highestHandEvaluator) getStats(ts model.TossSummary) model.TossStats {
	return ts.HandStats

}

func (*highestHandEvaluator) isBetter(old, new model.TossStats) bool {
	return willKeepHighestPotential(old, new)
}

// KeepHandHighestPotential will keep the hand with the highest potential score
func KeepHandHighestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, &highestHandEvaluator{})
}

var _ handEvaluator = (*lowestHandEvaluator)(nil)

type lowestHandEvaluator struct{}

func (*lowestHandEvaluator) getStats(ts model.TossSummary) model.TossStats {
	return ts.HandStats
}

func (*lowestHandEvaluator) isBetter(old, new model.TossStats) bool {
	return willKeepLowestPotential(old, new)
}

// KeepHandLowestPotential will keep the hand with the lowest potential score
func KeepHandLowestPotential(_ int, hand []model.Card) ([]model.Card, error) {
	return getEvaluatedHand(hand, &lowestHandEvaluator{})
}
