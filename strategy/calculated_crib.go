package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/scorer"
)

// GiveCribHighestPotential gives the crib the highest potential pointed crib
func GiveCribHighestPotential(_ int, hand []cards.Card) []cards.Card {
	isBetter := func(old, new float64) bool { return new > old }
	return getBestPotentialCrib(hand, isBetter)
}

// GiveCribLowestPotential gives the crib the lowest potential pointed hand
func GiveCribLowestPotential(_ int, hand []cards.Card) []cards.Card {
	isBetter := func(old, new float64) bool { return new < old }
	return getBestPotentialCrib(hand, isBetter)
}

func getBestPotentialCrib(hand []cards.Card, isBetter func(old, new float64) bool) []cards.Card {
	if len(hand) > 6 || len(hand) <= 4 {
		return nil
	}

	lenDeposit := len(hand) - 4
	bestCrib := make([]cards.Card, 0, lenDeposit)
	bestPotential := 0.0

	allDeposits := chooseFrom(lenDeposit, hand)

	seen := map[cards.Card]struct{}{}
	for _, c := range hand {
		seen[c] = struct{}{}
	}

	for i, dep := range allDeposits {
		p := getPotentialForDeposit(seen, dep)
		if i == 0 || isBetter(bestPotential, p) {
			bestCrib = bestCrib[:0]
			bestCrib = append(bestCrib, dep...)
			bestPotential = p
		}
	}

	return bestCrib
}

func getPotentialForDeposit(prevSeen map[cards.Card]struct{}, cribDeposit []cards.Card) float64 {
	seen := map[cards.Card]struct{}{}
	for k := range prevSeen {
		seen[k] = struct{}{}
	}
	for _, c := range cribDeposit {
		seen[c] = struct{}{}
	}

	totalCribPoints := 0
	totalHands := 0

	for i := 0; i < 52; i++ {
		lead := cards.NewCardFromNumber(i)
		if _, ok := seen[lead]; ok {
			continue
		}

		seen[lead] = struct{}{}

		options := otherOptions(4-len(cribDeposit), seen)

		for _, o := range options {
			totalCribPoints += scorer.CribPoints(lead, append(o, cribDeposit...))
			totalHands++
		}

		delete(seen, lead)
	}

	return float64(totalCribPoints) / float64(totalHands)
}
