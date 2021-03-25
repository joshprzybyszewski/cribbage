package play

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

func handContains(hand []model.Card, c model.Card) bool {
	for _, hc := range hand {
		if hc == c {
			return true
		}
	}
	return false
}

func hasBeenPegged(pegged []model.PeggedCard, c model.Card) bool {
	for _, pc := range pegged {
		if pc.Card == c {
			return true
		}
	}
	return false
}

func minUnpeggedValue(hand []model.Card, pegged []model.PeggedCard) int {
	peggedMap := make(map[model.Card]struct{}, len(pegged))
	for _, pc := range pegged {
		peggedMap[pc.Card] = struct{}{}
	}

	min := model.MaxPeggingValue + 1
	for _, hc := range hand {
		if _, ok := peggedMap[hc]; ok {
			continue
		}
		if pv := hc.PegValue(); pv < min {
			min = pv
		}
	}
	return min
}
