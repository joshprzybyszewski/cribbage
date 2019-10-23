package strategy

import (
	"sort"

	"github.com/joshprzybyszewski/cribbage/cards"
)

// PegToFifteen returns a card that yields a fifteen if it can
func PegToFifteen(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return pegToTarget(hand, prevPegs, curPeg, 15)
}

// PegToThirtyOne returns a card that yields 31 if it can
func PegToThirtyOne(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return pegToTarget(hand, prevPegs, curPeg, 31)
}

func pegToTarget(hand, prevPegs []cards.Card, curPeg, target int) (_ cards.Card, sayGo bool) {
	ct := 0
	for _, c := range hand {
		pegVal := c.PegValue() + curPeg
		switch {
		case pegVal == target:
			return c, false
		case pegVal > 31:
			ct++
		}
	}
	if ct == len(hand) {
		return cards.Card{}, true
	}
	return hand[0], false
}

// PegToPair returns a card from the hand iff that card makes a pair and does not push the count over 31
func PegToPair(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	ct := 0
	lastCard := prevPegs[len(prevPegs)-1]
	for _, c := range hand {
		if curPeg+c.PegValue() > 31 {
			ct++
		} else if c.Value == lastCard.Value {
			return c, false
		}
	}
	if ct == len(hand) {
		return cards.Card{}, true
	}
	return hand[0], false
}

// PegToRun returns a card that forms the longest run if one is possible
func PegToRun(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	peg := curPeg
	index := 0
	for i := len(prevPegs) - 1; i > 0; i-- {
		peg -= prevPegs[i].PegValue()
		if peg <= 0 {
			index = i
			break
		}
	}
	cardsToAnalyze := prevPegs[index:]
	for i := range cardsToAnalyze {
		cards := cardsToAnalyze[i:]
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Value < cards[j].Value
		})
	}
	return cards.Card{}, true
}
