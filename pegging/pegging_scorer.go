package pegging

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
)

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

	return points, nil
}

func validatePrevCards(prevCards []cards.Card, c cards.Card) error {
	existingCards := map[string]struct{}{}
	existingCards[c.String()] = struct{}{}
	for _, pc := range prevCards {
		if _, ok := existingCards[pc.String()]; ok {
			return errors.New(`it's impossible to peg the same card twice`)
		}
		existingCards[pc.String()] = struct{}{}
	}

	return nil
}
