package suggestions

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

type TossSummary struct {
	Tossed []model.Card

	HandStats *TossStats
	CribStats *TossStats
}

type TossStats struct {
	allPts []int

	min    int
	avg    float64
	median float64
	max    int
}

func (ts *TossStats) add(pts int) {
	ts.allPts = append(ts.allPts, pts)
}

func (ts *TossStats) calculate() {
	ts.min = 1000

	sum := 0

	for _, pt := range ts.allPts {
		sum += pt
		if pt < ts.min {
			ts.min = pt
		}
		if pt > ts.max {
			ts.max = pt
		}
	}
	ts.avg = float64(sum) / float64(len(ts.allPts))
}

func (ts *TossStats) Min() int {
	return ts.min
}
func (ts *TossStats) Median() float64 {
	return ts.median
}
func (ts *TossStats) Avg() float64 {
	return ts.avg
}
func (ts *TossStats) Max() int {
	return ts.max
}

func GetAllTosses(
	hand []model.Card,
) ([]TossSummary, error) {
	if len(hand) > 6 || len(hand) <= 4 {
		return nil, errors.New(`hand size must be between 4 and 6`)
	}

	allHands, err := chooseNFrom(4, hand)
	if err != nil {
		return nil, err
	}

	sums := []TossSummary{}

	for _, h := range allHands {
		tossed := without(hand, h)
		handStats, cribStats := getStatsForHand(h, tossed)

		sums = append(sums, TossSummary{
			Tossed:    tossed,
			HandStats: handStats,
			CribStats: cribStats,
		})
	}

	return sums, nil
}

func getStatsForHand(
	hand, tossed []model.Card,
) (handStats, cribStats *TossStats) {
	exclude := map[model.Card]struct{}{}
	for _, c := range hand {
		exclude[c] = struct{}{}
	}
	for _, c := range tossed {
		exclude[c] = struct{}{}
	}

	handStats = &TossStats{}
	cribStats = &TossStats{}
	defer handStats.calculate()
	defer cribStats.calculate()

	for i := 0; i < 52; i++ {
		lead := model.NewCardFromNumber(i)
		if _, ok := exclude[lead]; ok {
			continue
		}

		handStats.add(scorer.HandPoints(lead, hand))

		exclude[lead] = struct{}{}
		options := otherOptions(4-len(tossed), exclude)
		for _, o := range options {
			cribStats.add(scorer.CribPoints(lead, append(o, tossed...)))
		}
		delete(exclude, lead)
	}

	return handStats, cribStats
}
