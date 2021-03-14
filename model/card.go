package model

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
)

var (
	InvalidCard = Card{}
)

func NewCardFromString(card string) Card {
	c, err := NewCardFromExternalString(card)
	if err != nil {
		log.Printf(`dev error: NewCardFromString invalid card (%q): %s`, card, err.Error())
	}
	return c
}

// NewCardFromExternalString returns a card, or an error if the input is invalid
// Use this for external inputs (i.e. REST requests)
func NewCardFromExternalString(card string) (Card, error) {
	if len(card) > 3 {
		// Cards are expected to be of the form "AH" or "13C" or "5D".
		// Therfore, we don't support strings that have a len > 3.
		return InvalidCard, errors.New(`unknown card`)
	}

	value, err := getCardValue(card)
	if err != nil {
		return InvalidCard, err
	}

	suit, err := getSuit(card)
	if err != nil {
		return InvalidCard, err
	}

	return NewCard(suit, value), nil
}

func getCardValue(card string) (int, error) {
	switch string(card[0]) {
	case `A`, `a`:
		return 1, nil
	case `J`, `j`:
		return 11, nil
	case `Q`, `q`:
		return 12, nil
	case `K`, `k`:
		return 13, nil
	case `1`:
		// try parsing 10, 11, 12, or 13
		value, err := strconv.Atoi(card[:2])
		if err == nil {
			if value > 13 {
				return value, errors.New(`invalid card value`)
			}
			return value, nil
		}
		return 1, nil
	default:
		return strconv.Atoi(string(card[0]))
	}
}

func getSuit(card string) (Suit, error) {
	rs := []rune(card)
	suitStr := rs[len(rs)-1]
	switch string(suitStr) {
	case `S`, `s`:
		return Spades, nil
	case `C`, `c`:
		return Clubs, nil
	case `D`, `d`:
		return Diamonds, nil
	case `H`, `h`:
		return Hearts, nil
	default:
		return 0, errors.New(`bad input card: ` + card)
	}
}

func NewCardFromTinyInt(val int8) (Card, error) {
	return newCardFromNumber(int(val))
}

func NewCardFromNumber(val int) Card {
	c, err := newCardFromNumber(val)
	if err != nil {
		// set the card to the zero value (invalid card)
		c = InvalidCard
	}
	return c
}

func newCardFromNumber(val int) (Card, error) {
	if val < 0 || val > 51 {
		return InvalidCard, fmt.Errorf(`invalid num: %d`, val)
	}

	return NewCard(Suit(val/13), (val%13)+1), nil
}

func NewCard(suit Suit, value int) Card {
	return Card{
		Suit:  suit,
		Value: value,
	}
}

func (c Card) ToTinyInt() int8 {
	empty := Card{}
	if c == empty {
		return -1
	}
	return int8((int(c.Suit) * 13) + (c.Value - 1))
}

func (c Card) String() string {
	var val string
	switch c.Value {
	case 1:
		val = `A`
	case 11:
		val = `J`
	case 12:
		val = `Q`
	case 13:
		val = `K`
	default:
		val = strconv.Itoa(c.Value)
	}

	switch c.Suit {
	case Spades:
		val += `S`
	case Clubs:
		val += `C`
	case Diamonds:
		val += `D`
	case Hearts:
		val += `H`
	}

	return val
}

func (c Card) PegValue() int {
	if c.Value >= 10 {
		return 10
	}
	return c.Value
}

// SortByValue sorts a slice of cards either ascending or descending by their rank order
func SortByValue(input []Card, descending bool) []Card {
	retCards := make([]Card, len(input))
	_ = copy(retCards, input)
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
