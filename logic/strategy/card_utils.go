package strategy

import (
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

func countBits(num uint) (int, []int) {
	n := 0
	iter := 0
	idx := make([]int, 0)
	for num > 0 {
		if num&1 > 0 {
			n++
			idx = append(idx, iter)
		}
		num >>= 1
		iter++
	}
	return n, idx
}

func chooseFrom(desired int, hand []model.Card) [][]model.Card {
	if desired > len(hand) || desired > 4 || desired <= 0 {
		return nil
	}
	// the min int we need is the number with the lowest _n_ bits set, where n = desired
	// the max int we need is the number with the highest _n_ bits set, where n = desired
	// e.g. for 6 choose 4, we need min = 001111 and max = 111100
	hands := make([][]model.Card, 0)
	min := uint(1<<uint(desired)) - 1
	max := min << uint(len(hand)-desired)
	for i := min; i <= max; i++ {
		if n, idx := countBits(i); n == desired {
			thisHand := make([]model.Card, desired)
			for j, k := range idx {
				thisHand[j] = hand[k]
			}
			hands = append(hands, thisHand)
		}
	}
	return hands
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
