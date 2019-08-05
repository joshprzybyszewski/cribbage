package cards

import (
	"errors"
	"strconv"
)

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit      Suit
	Value     int
	deckValue int
}

func NewCardFromString(card string) Card {
	value := 0
	var err error

	suitStr := string(card[1:])
	switch string(card[0]) {
	case `A`:
		value = 1
	case `J`:
		value = 11
	case `Q`:
		value = 12
	case `K`:
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
		Suit:      suit,
		Value:     value,
		deckValue: (int(suit) * 13) + (value - 1),
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
