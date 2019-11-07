package cards

import (
	"sort"
)

// SortByValue sorts a slice of cards either ascending or descending by their rank order
func SortByValue(input []Card, descending bool) []Card {
	retCards := make([]Card, len(input))
	for i, c := range input {
		retCards[i] = c
	}
	sort.Slice(retCards, func(i, j int) bool {
		if retCards[i].Value == retCards[j].Value {
			if descending {
				return retCards[i].Suit > retCards[j].Suit
			}
			return retCards[i].Suit < retCards[j].Suit
		}
		if descending {
			return retCards[i].Value > retCards[j].Value
		}
		return retCards[i].Value < retCards[j].Value
	})
	return retCards
}
