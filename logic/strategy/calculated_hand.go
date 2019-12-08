package strategy

import (
	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

// KeepHandHighestPotential will keep the hand with the highest potential score
func KeepHandHighestPotential(_ int, hand []model.Card) []model.Card {
	isBetter := func(old, new float64) bool { return new > old }
	return getBestPotentialHand(hand, isBetter)
}

// KeepHandLowestPotential will keep the hand with the lowest potential score
func KeepHandLowestPotential(_ int, hand []model.Card) []model.Card {
	isBetter := func(old, new float64) bool { return new < old }
	return getBestPotentialHand(hand, isBetter)
}

func getBestPotentialHand(hand []model.Card, isBetter func(old, new float64) bool) []model.Card {
	if len(hand) > 6 || len(hand) <= 4 {
		return nil
	}

	bestHand := make([]model.Card, 0, 4)
	bestPotential := 0.0

	// TODO chooseFrom is producing hands which duplicate cards from the original
	// hand (e.g. with a hand of 2s, 3s, 4s, 5s, 6s, 7s, the first hand in allHands
	// was 2s, 3s, 3s, 3s)
	allHands := chooseFrom(4, hand)

	seen := map[model.Card]struct{}{}
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

func getHandPotentialForCribDeposit(prevSeen map[model.Card]struct{}, hand []model.Card) float64 {
	seen := map[model.Card]struct{}{}
	for k := range prevSeen {
		seen[k] = struct{}{}
	}
	for _, c := range hand {
		seen[c] = struct{}{}
	}

	totalHandPoints := 0
	totalHands := 0

	for i := 0; i < 52; i++ {
		lead := model.NewCardFromNumber(i)
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
