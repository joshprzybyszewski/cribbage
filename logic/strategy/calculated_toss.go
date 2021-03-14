package strategy

import (
	"github.com/joshprzybyszewski/cribbage/logic/suggestions"
	"github.com/joshprzybyszewski/cribbage/model"
)

type handEvaluator interface {
	getStats(model.TossSummary) model.TossStats
	isBetter(old, new model.TossStats) bool
}

func willKeepLowestPotential(old, new model.TossStats) bool {
	if old == nil {
		return true
	}
	if new == nil {
		return false
	}
	if new.Median() != old.Median() {
		return new.Median() < old.Median()
	}
	if new.Avg() != old.Avg() {
		return new.Avg() < old.Avg()
	}
	return new.Min() < old.Min()
}

func willKeepHighestPotential(old, new model.TossStats) bool {
	if old == nil {
		return true
	}
	if new == nil {
		return false
	}
	if new.Median() != old.Median() {
		return new.Median() > old.Median()
	}
	if new.Avg() != old.Avg() {
		return new.Avg() > old.Avg()
	}
	return new.Max() > old.Max()
}

func getEvaluatedHand(
	hand []model.Card,
	he handEvaluator,
) ([]model.Card, error) {
	summaries, err := suggestions.GetAllTosses(hand)
	if err != nil {
		return nil, err
	}

	lenDeposit := len(hand) - 4
	bestThrow := make([]model.Card, 0, lenDeposit)
	var prevBest model.TossStats

	for i := range summaries {
		stats := he.getStats(summaries[i])
		if he.isBetter(prevBest, stats) {
			bestThrow = bestThrow[:0]
			bestThrow = append(bestThrow, summaries[i].Tossed...)
			prevBest = stats
		}
	}

	return bestThrow, nil
}
