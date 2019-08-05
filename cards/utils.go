package cards

import (
	"sort"
)

func CardToInt8(card Card) int8 {
	return int8(card.deckValue)
}

func Int8ToCard(bits int8) Card {
	return NewCardFromNumber(int(bits))
}

func HandToInt32(hand []Card) int32 {
	if len(hand) != 4 {
		// only accepts 4 card hands
		return -1
	}

	deckValues := make([]int, 4)
	for i, c := range hand {
		deckValues[i] = c.deckValue
	}

	sort.Ints(deckValues)

	var bits int32

	bits = int32(deckValues[0]<<24 |
		deckValues[1]<<16 |
		deckValues[2]<<8 |
		deckValues[3]<<0)

	return bits
}

func Int32ToHand(bits int32) []Card {
	hand := make([]Card, 4)

	hand[0] = NewCardFromNumber(int((bits >> 24) & 0xFF))
	hand[1] = NewCardFromNumber(int((bits >> 16) & 0xFF))
	hand[2] = NewCardFromNumber(int((bits >> 8) & 0xFF))
	hand[3] = NewCardFromNumber(int((bits >> 0) & 0xFF))

	return hand
}
