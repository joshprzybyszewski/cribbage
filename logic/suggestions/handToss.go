package suggestions

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

func GetAllTosses(
	hand []model.Card,
) ([]model.TossSummary, error) {
	if len(hand) > 6 || len(hand) < 4 {
		return nil, errors.New(`hand size must be either 5 or 6`)
	}
	if containsDuplicates(hand) {
		return nil, errors.New(`hand contains duplicates`)
	}

	allHands, err := chooseNFrom(4, hand)
	if err != nil {
		return nil, err
	}

	summaries := []model.TossSummary{}

	for _, h := range allHands {
		tossed := without(hand, h)
		handStats, cribStats := getStatsForHand(h, tossed)

		summaries = append(summaries, model.TossSummary{
			Kept:      h,
			Tossed:    tossed,
			HandStats: handStats,
			CribStats: cribStats,
		})
	}

	return summaries, nil
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
		cutCard := model.NewCardFromNumber(i)
		if _, ok := exclude[cutCard]; ok {
			continue
		}

		handStats.add(scorer.HandPoints(cutCard, hand))

		exclude[cutCard] = struct{}{}
		options := otherOptions(4-len(tossed), exclude)
		for _, o := range options {
			cribStats.add(scorer.CribPoints(cutCard, append(o, tossed...)))
		}
		delete(exclude, cutCard)
	}

	return handStats, cribStats
}

func containsDuplicates(hand []model.Card) bool {
	found := map[model.Card]struct{}{}
	for _, c := range hand {
		if _, ok := found[c]; ok {
			return true
		}
		found[c] = struct{}{}
	}
	return false
}
