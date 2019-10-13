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
	if totalPegged+c.Value > 31 {
		return 0, nil
	}
	cardsToAnalyze := prevCards[indexOfCardsToUse:]

	points := 0
	runPoints, err := scoreRun(cardsToAnalyze, c)
	if err != nil {
		return 0, err
	}
	points += runPoints

	pairPoints := scorePairs(cardsToAnalyze, c)
	points += pairPoints

	switch totalPegged + c.PegValue() {
	case 15, 31:
		points += 2
	}

	return points, nil
}

func scoreRun(cardsToAnalyze []cards.Card, c cards.Card) (int, error) {
	runLen := 0
	for i := len(cardsToAnalyze) - 2; i >= 0; i-- {
		if !isRun(append(cardsToAnalyze[i:], c)) {
			continue
		}
		runLen = (len(cardsToAnalyze) - i) + 1
	}
	if runLen >= 8 {
		return 0, errRunTooLong
	} else if runLen >= 3 {
		return runLen, nil
	}
	return 0, nil
}

func isRun(c []cards.Card) bool {
	sortedCards := make([]cards.Card, 0, len(c))
	for _, card := range c {
		sortedCards = append(sortedCards, cards.NewCardFromString(card.String()))
	}
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Value > sortedCards[j].Value
	})
	for i := 1; i < len(sortedCards); i++ {
		if sortedCards[i-1].Value != sortedCards[i].Value+1 {
			return false
		}
	}
	return true
}

func scorePairs(prevCards []cards.Card, c cards.Card) int {
	points := 0
	for i := len(prevCards) - 1; i >= 0; i-- {
		if prevCards[i].Value != c.Value {
			break
		}
		points += 2 * (len(prevCards) - i)
	}
	return points
}

func validatePrevCards(prevCards []cards.Card, c cards.Card) error {
	if len(prevCards) >= 4*4 {
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
