package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/scorer"
)

// KeepHandHighestPotential will keep the hand with the highest potential score
func KeepHandHighestPotential(_ int, hand []cards.Card) []cards.Card {
	isBetter := func(old, new float64) bool { return new > old }
	return getBestPotentialHand(hand, isBetter)
}

// KeepHandLowestPotential will keep the hand with the lowest potential score
func KeepHandLowestPotential(_ int, hand []cards.Card) []cards.Card {
	isBetter := func(old, new float64) bool { return new < old }
	return getBestPotentialHand(hand, isBetter)
}

func getBestPotentialHand(hand []cards.Card, isBetter func(old, new float64) bool) []cards.Card {
	if len(hand) > 6 || len(hand) <= 4 {
		return nil
	}

	bestHand := make([]cards.Card, 0, 4)
	bestPotential := 0.0

	allHands := chooseFrom(4, hand)

	seen := map[cards.Card]struct{}{}
	for _, c := range hand {
		seen[c] = struct{}{}
	}

	for i, h := range allHands {
		p := getHandPotentialForCribDeposit(seen, h)
		if i == 0 || isBetter(bestPotential, p) {
			bestHand = bestHand[:0]
			bestHand = append(bestHand, h...)
			bestPotential = p
		}
	}

	return without(hand, bestHand)
}

func getHandPotentialForCribDeposit(prevSeen map[cards.Card]struct{}, hand []cards.Card) float64 {
	seen := map[cards.Card]struct{}{}
	for k := range prevSeen {
		seen[k] = struct{}{}
	}
	for _, c := range hand {
		seen[c] = struct{}{}
	}

	totalHandPoints := 0
	totalHands := 0

	for i := 0; i < 52; i++ {
		lead := cards.NewCardFromNumber(i)
		if _, ok := seen[lead]; ok {
			continue
		}

		seen[lead] = struct{}{}

		totalHandPoints += scorer.HandPoints(lead, hand)
		totalHands++

		delete(seen, lead)
	}

	return float64(totalHandPoints) / float64(totalHands)
}
