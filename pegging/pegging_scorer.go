package pegging

import (
	"errors"
	"sort"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var (
	errSameCardTwice = errors.New(`it's impossible to peg the same card twice`)
	errTooManyCards  = errors.New(`it's impossible to have this many cards pegged ever`)
)

// PointsForCard needs a real comment to make the linter happier
func PointsForCard(prevCards []cards.Card, c cards.Card) (int, error) {
	if err := validatePrevCards(prevCards, c); err != nil {
		return 0, err
	}

	totalPegged := 0
	indexOfCardsToUse := 0
	for i, pc := range prevCards {
		totalPegged += pc.PegValue()
		if totalPegged > 31 {
			totalPegged = pc.PegValue()
			indexOfCardsToUse = i
		}
	}

	points := 0
	switch totalPegged + c.PegValue() {
	case 15, 31:
		points += 2
	}

	cardsToAnalyze := prevCards[indexOfCardsToUse:]
	for i := len(cardsToAnalyze) - 1; i >= 0; i-- {
		if cardsToAnalyze[i].Value != c.Value {
			break
		}
		points += 2 * (len(cardsToAnalyze) - i)
	}
	if points > 0 {
		return points, nil
	}
	sortedCards := append(cardsToAnalyze, c)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Value < sortedCards[j].Value
	})
	points = 1
	for i := len(sortedCards) - 1; i >= 1; i-- {
		if sortedCards[i].Value == sortedCards[i-1].Value+1 {
			points++
		} else {
			break
		}
	}
	if points == 1 {
		return 0, nil
	}
	return points, nil
}

func validatePrevCards(prevCards []cards.Card, c cards.Card) error {
	if len(prevCards) > 4*4 {
		// 4 players can each peg four cards. that's our max
		return errTooManyCards
	}

	existingCards := map[string]struct{}{}
	existingCards[c.String()] = struct{}{}
	for _, pc := range prevCards {
		if _, ok := existingCards[pc.String()]; ok {
			return errSameCardTwice
		}
		existingCards[pc.String()] = struct{}{}
	}

	return nil
}
