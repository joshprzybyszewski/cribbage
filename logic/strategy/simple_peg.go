package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

// PegToFifteen returns a card that yields a fifteen if it can
func PegToFifteen(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (_ model.Card, sayGo bool) {
	return pegToTarget(hand, prevPegs, curPeg, 15)
}

// PegToThirtyOne returns a card that yields 31 if it can
func PegToThirtyOne(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (_ model.Card, sayGo bool) {
	return pegToTarget(hand, prevPegs, curPeg, 31)
}

func pegToTarget(hand []model.Card, prevPegs []model.PeggedCard, curPeg, target int) (_ model.Card, sayGo bool) {
	if mustSayGo(hand, curPeg) {
		return model.Card{}, true
	}
	for _, c := range hand {
		if c.PegValue()+curPeg == target {
			return c, false
		}
	}
	return firstValidCard(hand, curPeg)
}

func firstValidCard(hand []model.Card, curPeg int) (model.Card, bool) {
	for _, c := range hand {
		if c.PegValue()+curPeg <= 31 {
			return c, false
		}
	}
	return model.Card{}, true
}

// PegToPair returns a card from the hand iff that card makes a pair and does not push the count over 31
func PegToPair(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (_ model.Card, sayGo bool) {
	if mustSayGo(hand, curPeg) {
		return model.Card{}, true
	}
	lastCard := prevPegs[len(prevPegs)-1]
	for _, c := range hand {
		if c.Value == lastCard.Value && curPeg+c.PegValue() <= 31 {
			return c, false
		}
	}
	return firstValidCard(hand, curPeg)
}

// PegToRun returns a card that forms the longest run if one is possible
func PegToRun(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (_ model.Card, sayGo bool) {
	if mustSayGo(hand, curPeg) {
		return model.Card{}, true
	}
	// Runs reset at 31 or go, so only look at the cards since one of those have happened
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
	if len(cardsToAnalyze) < 2 {
		return firstValidCard(hand, curPeg)
	}
	for i := range cardsToAnalyze {
		for _, c := range hand {
			handCopy := make([]model.Card, 0, len(cardsToAnalyze)-i+1)
			handCopy = append(handCopy, c)
			for _, pc := range cardsToAnalyze[i:] {
				handCopy = append(handCopy, pc.Card)
			}
			runCards := model.SortByValue(handCopy, false)
			for j := 0; j < len(runCards)-1; j++ {
				if runCards[j].Value != runCards[j+1].Value-1 {
					break
				}
				if j == len(runCards)-2 {
					return c, false
				}
			}
		}
	}
	return hand[0], false
}

func mustSayGo(hand []model.Card, curPeg int) bool {
	ct := 0
	for _, c := range hand {
		if curPeg+c.PegValue() > 31 {
			ct++
		}
	}
	return ct == len(hand)
}
