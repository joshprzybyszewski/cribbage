package model

import (
	"errors"
	"sort"
	"strconv"
)

func NewCardFromString(card string) Card {
	value := 0
	var err error

	suitStr := string(card[1:])
	switch string(card[0]) {
	case `A`, `a`:
		value = 1
	case `J`, `j`:
		value = 11
	case `Q`, `q`:
		value = 12
	case `K`, `k`:
		value = 13
	case `1`:
		// try parsing 10, 11, 12, or 13
		value, err = strconv.Atoi(string(card[:2]))
		if err == nil {
			suitStr = string(card[2:])
		} else {
			value, err = strconv.Atoi(string(card[0]))
		}
	default:
		value, err = strconv.Atoi(string(card[0]))
	}

	var suit Suit
	switch suitStr {
	case `S`, `s`, `♤`, `♠︎`:
		suit = Spades
	case `C`, `c`, `♧`, `♣︎`:
		suit = Clubs
	case `D`, `d`, `♢`, `♦`:
		suit = Diamonds
	case `H`, `h`, `♡`, `♥︎`:
		suit = Hearts
	default:
		err = errors.New(`bad input card: ` + card)
	}

	if err != nil {
		println(`got an error! ` + err.Error())
		return Card{}
	}

	return NewCard(suit, value)
}

func NewCardFromNumber(val int) Card {
	if val < 0 || val > 51 {
		println(`cannot support this val`)
		return Card{}
	}

	return NewCard(Suit(val/13), (val%13)+1)
}

func NewCard(suit Suit, value int) Card {
	return Card{
		Suit:  suit,
		Value: value,
	}
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
		val += `♠︎`
	case Clubs:
		val += `♣︎`
	case Diamonds:
		val += `♦`
	case Hearts:
		val += `♥︎`
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
