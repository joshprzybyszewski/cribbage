package strategy

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

func otherOptions(desired int, avoid map[model.Card]struct{}) [][]model.Card {
	options := [][]model.Card{}

	for o1 := 0; o1 < 52; o1++ {
		oc1 := model.NewCardFromNumber(o1)
		if _, ok := avoid[oc1]; ok {
			continue
		}
		for o2 := o1 + 1; o2 < 52; o2++ {
			oc2 := model.NewCardFromNumber(o2)
			if _, ok := avoid[oc2]; ok {
				continue
			}

			if desired == 2 {
				options = append(options, []model.Card{oc1, oc2})
				continue
			}

			for o3 := o2 + 1; o3 < 52; o3++ {
				oc3 := model.NewCardFromNumber(o3)
				if _, ok := avoid[oc3]; ok {
					continue
				}
				options = append(options, []model.Card{oc1, oc2, oc3})
			}
		}
	}

	return options
}

func chooseFrom(k int, hand []model.Card) ([][]model.Card, error) {
	if len(hand) > 6 || k < 1 || k > len(hand) {
		return nil, errors.New(`invalid input`)
	}
	if k == 1 {
		all := make([][]model.Card, len(hand))
		for i, e := range hand {
			all[i] = []model.Card{e}
		}
		return all, nil
	}
	if k == len(hand) {
		cpy := make([]model.Card, len(hand))
		copy(cpy, hand)
		return [][]model.Card{cpy}, nil
	}
	// 6 choose 3 = 20, the max number of combos we would ever have
	all := make([][]model.Card, 0, 20)
	// for the first n-k cards, recursively find combinations of length k-1 which are
	// combined with the current card to get combinations of length k
	for i := 0; i <= len(hand)-k; i++ {
		c := hand[i]
		others := hand[i+1:]
		otherSets, err := chooseFrom(k-1, others)
		if err != nil {
			return nil, err
		}
		for _, s := range otherSets {
			set := make([]model.Card, 1, k)
			set[0] = c
			set = append(set, s...)
			all = append(all, set)
		}
	}
	return all, nil
}

// without returns the cards in superset minus the subsetToRemove
func without(superset, subsetToRemove []model.Card) []model.Card {
	removed := make([]model.Card, 0, len(superset)-len(subsetToRemove))
	rem := make(map[model.Card]struct{}, len(subsetToRemove))
	for _, s := range subsetToRemove {
		rem[s] = struct{}{}
	}
	for _, c := range superset {
		if _, ok := rem[c]; ok {
			continue
		}
		removed = append(removed, c)
	}
	return removed
}
