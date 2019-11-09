package model

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit  Suit
	Value int // Ace is 1, King is 13
}

const NumCardsPerDeck = 52

type PlayerID int64
type PlayerColor string

type Player struct {
	ID   PlayerID
	Name string
}

type GamePlayer struct {
	ID    PlayerID
	Color PlayerColor
}

type Phase uint8

const (
	Deal Phase = iota
	BuildCrib
	Cut
	Pegging
	Counting
	CribCounting
	Done
)

type Game struct {
	Players       []Player
	CurrentDealer PlayerID
	CurrentScores map[PlayerColor]uint8
	LagScores     map[PlayerColor]uint8
	Phase         Phase
	Hands         map[PlayerID][]Card
	Crib          []Card
}
