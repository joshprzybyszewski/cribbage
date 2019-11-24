package model

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit Suit `protobuf:"varint,1,req,name=suit,proto3" json:"suit"`
	// Ace is 1, King is 13
	Value int `protobuf:"varint,2,req,name=value,proto3" json:"value"`
}

const NumCardsPerDeck = 52
const JackValue = 11 // Ace is 1, King is 13

type PeggedCard struct {
	Card
	PlayerID PlayerID `protobuf:"varint,3,req,name=playerID,proto3" json:"playerID"`
}

type PlayerID uint32
type GameID uint32

const (
	InvalidPlayerID PlayerID = 0
	InvalidGameID   GameID   = 0
)

type PlayerColor int8

const (
	Green PlayerColor = iota
	Blue
	Red
)

func (c PlayerColor) String() string {
	switch c {
	case Blue:
		return `blue`
	case Red:
		return `red`
	case Green:
		return `green`
	}
	return `notacolor`
}

type Player struct {
	ID    PlayerID               `protobuf:"varint,1,req,name=id,proto3" json:"id"`
	Name  string                 `protobuf:"string,2,req,name=name,proto3" json:"name"`
	Games map[GameID]PlayerColor `protobuf:"map<varint, varint>,3,opt,name=games,proto3" json:"games,omitempty"`
}

type InteractionMeans struct {
	Means string      `protobuf:"-" json:"-"`
	Info  interface{} `protobuf:"-" json:"-"`
}

type Blocker int

const (
	DealCards Blocker = iota
	CribCard
	CutCard
	PegCard
	CountHand
	CountCrib
)

type CribBlocker struct {
	Desired      int                      `protobuf:"varint,1,req,name=desired,proto3" json:"desired"`
	Dealer       PlayerID                 `protobuf:"varint,2,req,name=playerID,proto3" json:"playerID"`
	PlayerColors map[PlayerID]PlayerColor `protobuf:"map<varint,varint>,3,opt,name=playerColors,proto3" json:"playerColors,omitempty"`
}

type PlayerAction struct {
	GameID    GameID      `protobuf:"varint,1,req,name=gameID,proto3" json:"gameID"`
	ID        PlayerID    `protobuf:"varint,2,req,name=id,proto3" json:"id"`
	Overcomes Blocker     `protobuf:"varint,3,req,name=overcomes,proto3" json:"overcomes"`
	Action    interface{} `protobuf:"-" json:"action"`
}

type DealAction struct {
	NumShuffles int `protobuf:"varint,1,req,name=numShuffles,proto3" json:"numShuffles"`
}

type BuildCribAction struct {
	Cards []Card `protobuf:"Card,rep,1,req,name=cards,proto3" json:"cards"`
}

type CutDeckAction struct {
	Percentage float64 `protobuf:"varint,1,req,name=percentage,proto3" json:"percentage"`
}

type PegAction struct {
	Card  Card `protobuf:"Card,1,req,name=card,proto3" json:"card"`
	SayGo bool `protobuf:"bool,2,req,name=sayGo,proto3" json:"sayGo"`
}

type CountHandAction struct {
	Pts int `protobuf:"varint,1,req,name=pts,proto3" json:"pts"`
}

type CountCribAction struct {
	Pts int `protobuf:"varint,1,req,name=pts,proto3" json:"pts"`
}

type Phase int

const (
	Deal Phase = iota
	BuildCribReady
	BuildCrib
	CutReady
	Cut
	PeggingReady
	Pegging
	CountingReady
	Counting
	CribCountingReady
	CribCounting
	DealingReady
)

const (
	WinningScore    int = 121
	MaxPeggingValue int = 31
)

const (
	MinPlayerGame int = 2
	MaxPlayerGame int = 4
)

type Game struct {
	ID              GameID                   `protobuf:"varint,1,name=id,proto3" json:"id"`
	NumActions      int                      `protobuf:"varint,2,name=numActions,proto3" json:"numActions"`
	Players         []Player                 `protobuf:"bytes,3,rep,name=players,proto3,proto3" json:"players"`
	Deck            Deck                     `protobuf:"-" json:"-"`
	BlockingPlayers map[PlayerID]Blocker     `protobuf:"map<varint, varint>,4,opt,name=blockingPlayers,proto3" json:"blockingPlayers,omitempty"`
	CurrentDealer   PlayerID                 `protobuf:"varint,5,opt,name=currentDealer,proto3" json:"currentDealer"`
	PlayerColors    map[PlayerID]PlayerColor `protobuf:"map<varint, varint>,6,name=playerColors,proto3" json:"playerColors,omitempty"`
	CurrentScores   map[PlayerColor]int      `protobuf:"map<varint, varint>,7,opt,name=currentScores,proto3" json:"currentScores"`
	LagScores       map[PlayerColor]int      `protobuf:"map<varint, varint>,8,opt,name=lagScores,proto3" json:"lagScores"`
	Phase           Phase                    `protobuf:"varint,9,opt,name=phase,proto3" json:"phase"`
	Hands           map[PlayerID][]Card      `protobuf:"map<varint, bytes>,10,opt,name=hands,proto3" json:"hands,omitempty"`
	CutCard         Card                     `protobuf:"Card,11,opt,name=cutCard,proto3" json:"cutCard,omitempty"`
	Crib            []Card                   `protobuf:"Card,12,rep,name=crib,proto3" json:"crib,omitempty"`
	PeggedCards     []PeggedCard             `protobuf:"PeggedCard,13,rep,opt,name=peggedCards,proto3" json:"peggedCards,omitempty"`
}
