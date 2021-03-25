package suggestions

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

func otherOptions(desired int, exclude map[model.Card]struct{}) [][]model.Card {
	options := [][]model.Card{}

	for o1 := 0; o1 < 52; o1++ {
		oc1 := model.NewCardFromNumber(o1)
		if _, ok := exclude[oc1]; ok {
			continue
		}
		for o2 := o1 + 1; o2 < 52; o2++ {
			oc2 := model.NewCardFromNumber(o2)
			if _, ok := exclude[oc2]; ok {
				continue
			}

			if desired == 2 {
				options = append(options, []model.Card{oc1, oc2})
				continue
			}

			for o3 := o2 + 1; o3 < 52; o3++ {
				oc3 := model.NewCardFromNumber(o3)
				if _, ok := exclude[oc3]; ok {
					continue
				}
				options = append(options, []model.Card{oc1, oc2, oc3})
			}
		}
	}

	return options
}

func chooseNFrom(n int, hand []model.Card) ([][]model.Card, error) {
	if n < 1 || n > len(hand) {
		return nil, errors.New(`developer error: invalid n`)
	}
	if len(hand) > 6 {
		return nil, errors.New(`too many cards in hand (maximum 6)`)
	}
	if n == 1 {
		all := make([][]model.Card, len(hand))
		for i, e := range hand {
			all[i] = []model.Card{e}
		}
		return all, nil
	}
	if n == len(hand) {
		cpy := make([]model.Card, len(hand))
		copy(cpy, hand)
		return [][]model.Card{cpy}, nil
	}
	// 6 choose 3 = 20, the max number of combos we would ever have
	all := make([][]model.Card, 0, 20)
	// for the first len(hand)-n cards, recursively find combinations of length n-1 which are
	// combined with the current card to get combinations of length n
	for i := 0; i <= len(hand)-n; i++ {
		c := hand[i]
		others := hand[i+1:]
		otherSets, err := chooseNFrom(n-1, others)
		if err != nil {
			return nil, err
		}
		for _, s := range otherSets {
			set := make([]model.Card, 1, n)
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
