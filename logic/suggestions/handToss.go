package suggestions

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

func GetAllTosses(
	hand []model.Card,
) ([]model.TossSummary, error) {
	if len(hand) > 6 || len(hand) <= 4 {
		return nil, errors.New(`hand size must be between 4 and 6`)
	}

	allHands, err := chooseNFrom(4, hand)
	if err != nil {
		return nil, err
	}

	sums := []model.TossSummary{}

	for _, h := range allHands {
		tossed := without(hand, h)
		handStats, cribStats := getStatsForHand(h, tossed)

		sums = append(sums, model.TossSummary{
			Kept:      h,
			Tossed:    tossed,
			HandStats: handStats,
			CribStats: cribStats,
		})
	}

	return sums, nil
}

func getStatsForHand(
	hand, tossed []model.Card,
) (handStats, cribStats *tossStats) {
	exclude := map[model.Card]struct{}{}
	for _, c := range hand {
		exclude[c] = struct{}{}
	}
	for _, c := range tossed {
		exclude[c] = struct{}{}
	}

	handStats = &tossStats{}
	cribStats = &tossStats{}
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
