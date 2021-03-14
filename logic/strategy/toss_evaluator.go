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
	if old == nil {
		return true
	}
	if new == nil {
		return false
	}

	switch te.kind {
	case highestIsBetter:
		return higherStats(old, new)

	case lowestIsBetter:
		return lowerStats(old, new)
	}
	return false
}

func higherStats(old, new model.TossStats) bool {
	if differentFloats(new.Median(), old.Median()) {
		return new.Median() > old.Median()
	}
	if differentFloats(new.Avg(), old.Avg()) {
		return new.Avg() > old.Avg()
	}
	if new.Max() != old.Max() {
		return new.Max() > old.Max()
	}
	return new.Min() > old.Min()
}

func lowerStats(old, new model.TossStats) bool {
	if differentFloats(new.Median(), old.Median()) {
		return new.Median() < old.Median()
	}
	if differentFloats(new.Avg(), old.Avg()) {
		return new.Avg() < old.Avg()
	}
	if new.Min() != old.Min() {
		return new.Min() < old.Min()
	}
	return new.Max() < old.Max()
}

func differentFloats(
	a, b float64,
) bool {
	epsilon := 0.001
	return a-b > epsilon || b-a > epsilon
}
