package strategy

import "github.com/joshprzybyszewski/cribbage/model"

type betterKind int

const (
	highestIsBetter betterKind = 0
	lowestIsBetter  betterKind = 1
)

var _ handEvaluator = (*tossEvaluator)(nil)

type tossEvaluator struct {
	forHand bool
	kind    betterKind
}

func newTossEvaluator(
	forHand bool,
	kind betterKind,
) handEvaluator {
	return &tossEvaluator{
		forHand: forHand,
		kind:    kind,
	}
}

func (te *tossEvaluator) getStats(ts model.TossSummary) model.TossStats {
	if te.forHand {
		return ts.HandStats
	}
	return ts.CribStats

}

func (te *tossEvaluator) isBetter(old, new model.TossStats) bool {
	switch te.kind {
	case highestIsBetter:
		return willKeepHighestPotential(old, new)
	case lowestIsBetter:
		return willKeepLowestPotential(old, new)
	}
	return false
}
