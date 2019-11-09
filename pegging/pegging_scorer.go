package pegging

import (
	"errors"
	"sort"

	"github.com/joshprzybyszewski/cribbage/model"
)

var (
	errSameCardTwice = errors.New(`it's impossible to peg the same card twice`)
	errTooManyCards  = errors.New(`it's impossible to have this many cards pegged ever`)
)

// PointsForCard returns how many points are received for the given card, provided the previously pegged cards
func PointsForCard(prevCards []model.Card, c model.Card) (int, error) {
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
	if totalPegged+c.PegValue() > 31 {
		// If this card pushes us over 31, then don't consider any previous cards
		cardsToAnalyze = cardsToAnalyze[:0]
	}

	points := 0
	runPoints := scoreRun(cardsToAnalyze, c)
	points += runPoints

	pairPoints := scorePairs(cardsToAnalyze, c)
	points += pairPoints

	switch totalPegged + c.PegValue() {
	case 15, 31:
		points += 2
	}

	return points, nil
}

func scoreRun(cardsToAnalyze []model.Card, c model.Card) int {
	runLen := 0
	for i := len(cardsToAnalyze) - 2; i >= 0; i-- {
		if !isRun(append(cardsToAnalyze[i:], c)) {
			continue
		}
		runLen = (len(cardsToAnalyze) - i) + 1
	}
	if runLen >= 3 {
		return runLen
	}
	return 0
}

func isRun(c []model.Card) bool {
	sortedCards := make([]model.Card, 0, len(c))
	for _, card := range c {
		sortedCards = append(sortedCards, model.NewCardFromString(card.String()))
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

func scorePairs(prevCards []model.Card, c model.Card) int {
	points := 0
	for i := len(prevCards) - 1; i >= 0; i-- {
		if prevCards[i].Value != c.Value {
			break
		}
		// this will add the correct number of points for pairs
		// the first time 2, then 4, then 6
		points += 2 * (len(prevCards) - i)
	}
	return points
}

func validatePrevCards(prevCards []model.Card, c model.Card) error {
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
