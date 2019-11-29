package model

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit Suit `protobuf:"varint,1,req,name=suit,proto3" json:"suit"` //nolint:lll
	// Ace is 1, King is 13
	Value int `protobuf:"varint,2,req,name=value,proto3" json:"value"` //nolint:lll
}

const NumCardsPerDeck = 52
const JackValue = 11 // Ace is 1, King is 13

type PeggedCard struct {
	Card
	Action   int      `protobuf:"varint,3,req,name=action,proto3" json:"action"`     //nolint:lll
	PlayerID PlayerID `protobuf:"varint,4,req,name=playerID,proto3" json:"playerID"` //nolint:lll
}

type PlayerID string
type GameID uint32

const (
	InvalidPlayerID PlayerID = ``
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
	ID    PlayerID               `protobuf:"varint,1,req,name=id,proto3" json:"id"`                              //nolint:lll
	Name  string                 `protobuf:"string,2,req,name=name,proto3" json:"name"`                          //nolint:lll
	Games map[GameID]PlayerColor `protobuf:"map<varint, varint>,3,opt,name=games,proto3" json:"games,omitempty"` //nolint:lll
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
	Desired      int                      `protobuf:"varint,1,req,name=desired,proto3" json:"desired"`                                 //nolint:lll//nolint:lll
	Dealer       PlayerID                 `protobuf:"varint,2,req,name=playerID,proto3" json:"playerID"`                               //nolint:lll//nolint:lll
	PlayerColors map[PlayerID]PlayerColor `protobuf:"map<varint,varint>,3,opt,name=playerColors,proto3" json:"playerColors,omitempty"` //nolint:lll//nolint:lll
}

type PlayerAction struct {
	GameID    GameID      `protobuf:"varint,1,req,name=gameID,proto3" json:"gameID"`       //nolint:lll
	ID        PlayerID    `protobuf:"varint,2,req,name=id,proto3" json:"id"`               //nolint:lll
	Overcomes Blocker     `protobuf:"varint,3,req,name=overcomes,proto3" json:"overcomes"` //nolint:lll
	Action    interface{} `protobuf:"-" json:"action"`                                     //nolint:lll
}

type DealAction struct {
	NumShuffles int `protobuf:"varint,1,req,name=numShuffles,proto3" json:"numShuffles"` //nolint:lll
}

type BuildCribAction struct {
	Cards []Card `protobuf:"Card,rep,1,req,name=cards,proto3" json:"cards"` //nolint:lll
}

type CutDeckAction struct {
	Percentage float64 `protobuf:"varint,1,req,name=percentage,proto3" json:"percentage"` //nolint:lll
}

type PegAction struct {
	Card  Card `protobuf:"Card,1,req,name=card,proto3" json:"card"`   //nolint:lll
	SayGo bool `protobuf:"bool,2,req,name=sayGo,proto3" json:"sayGo"` //nolint:lll
}

type CountHandAction struct {
	Pts int `protobuf:"varint,1,req,name=pts,proto3" json:"pts"` //nolint:lll
}

type CountCribAction struct {
	Pts int `protobuf:"varint,1,req,name=pts,proto3" json:"pts"` //nolint:lll
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
	ID              GameID                   `protobuf:"varint,1,name=id,proto3" json:"id"`                                                      //nolint:lll
	Players         []Player                 `protobuf:"bytes,3,rep,name=players,proto3,proto3" json:"players"`                                  //nolint:lll
	Deck            Deck                     `protobuf:"-" json:"-"`                                                                             //nolint:lll
	BlockingPlayers map[PlayerID]Blocker     `protobuf:"map<varint, varint>,4,opt,name=blockingPlayers,proto3" json:"blockingPlayers,omitempty"` //nolint:lll
	CurrentDealer   PlayerID                 `protobuf:"varint,5,opt,name=currentDealer,proto3" json:"currentDealer"`                            //nolint:lll
	PlayerColors    map[PlayerID]PlayerColor `protobuf:"map<varint, varint>,6,name=playerColors,proto3" json:"playerColors,omitempty"`           //nolint:lll
	CurrentScores   map[PlayerColor]int      `protobuf:"map<varint, varint>,7,opt,name=currentScores,proto3" json:"currentScores"`               //nolint:lll
	LagScores       map[PlayerColor]int      `protobuf:"map<varint, varint>,8,opt,name=lagScores,proto3" json:"lagScores"`                       //nolint:lll
	Phase           Phase                    `protobuf:"varint,9,opt,name=phase,proto3" json:"phase"`                                            //nolint:lll
	Hands           map[PlayerID][]Card      `protobuf:"map<varint, bytes>,10,opt,name=hands,proto3" json:"hands,omitempty"`                     //nolint:lll
	CutCard         Card                     `protobuf:"Card,11,opt,name=cutCard,proto3" json:"cutCard,omitempty"`                               //nolint:lll
	Crib            []Card                   `protobuf:"Card,12,rep,name=crib,proto3" json:"crib,omitempty"`                                     //nolint:lll
	PeggedCards     []PeggedCard             `protobuf:"PeggedCard,13,rep,opt,name=peggedCards,proto3" json:"peggedCards,omitempty"`             //nolint:lll

	actions []PlayerAction
}
