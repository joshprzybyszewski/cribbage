package pegging

import (
	"errors"
	"sort"

	"github.com/joshprzybyszewski/cribbage/cards"
)

var (
	errSameCardTwice = errors.New(`it's impossible to peg the same card twice`)
	errTooManyCards  = errors.New(`it's impossible to have this many cards pegged ever`)
	errRunTooLong    = errors.New(`we cannot have a run of 8 in pegging`)
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
	cardsToAnalyze := prevCards[indexOfCardsToUse:]

	points := 0
	runPoints, err := scoreRun(cardsToAnalyze, c)
	if err != nil {
		return 0, err
	}
	points += runPoints
	switch totalPegged + c.PegValue() {
	case 15, 31:
		points += 2
	}

	for i := len(cardsToAnalyze) - 1; i >= 0; i-- {
		if cardsToAnalyze[i].Value != c.Value {
			break
		}
		points += 2 * (len(cardsToAnalyze) - i)
	}
	return points, nil
}

func scoreRun(cardsToAnalyze []cards.Card, c cards.Card) (int, error) { // Check for runs
	runLen := 0
	for i := len(cardsToAnalyze) - 1; i >= 0; i-- {
		// Make a slice of the first i+1 cards in new memory
		sortedCards := make([]cards.Card, 0, i+1)
		sortedCards = append(sortedCards, cardsToAnalyze[i:]...)
		sortedCards = append(sortedCards, c)
		sort.Slice(sortedCards, func(i, j int) bool {
			return sortedCards[i].Value > sortedCards[j].Value
		})
		for j := 1; j < len(sortedCards); j++ {
			if sortedCards[j-1].Value == sortedCards[j].Value+1 {
				runLen = j + 1
			} else {
				break
			}
		}
	}
	if runLen >= 8 {
		return 0, errRunTooLong
	} else if runLen >= 3 {
		return runLen, nil
	}
	return 0, nil
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
