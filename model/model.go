package model

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit Suit `protobuf:"varint,1,req,name=suit,proto3" json:"suit" bson:"suit"` //nolint:lll
	// Ace is 1, King is 13
	Value int `protobuf:"varint,2,req,name=value,proto3" json:"value" bson:"value"` //nolint:lll
}

const NumCardsPerDeck = 52
const JackValue = 11 // Ace is 1, King is 13

type PeggedCard struct {
	Card
	Action   int      `protobuf:"varint,3,req,name=action,proto3" json:"action" bson:"action"`       //nolint:lll
	PlayerID PlayerID `protobuf:"varint,4,req,name=playerID,proto3" json:"playerID" bson:"playerID"` //nolint:lll
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
	ID    PlayerID               `protobuf:"varint,1,req,name=id,proto3" json:"id" bson:"id"`                                 //nolint:lll
	Name  string                 `protobuf:"string,2,req,name=name,proto3" json:"name" bson:"name"`                           //nolint:lll
	Games map[GameID]PlayerColor `protobuf:"map<varint, varint>,3,opt,name=games,proto3" json:"games,omitempty" bson:"games"` //nolint:lll
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
	Desired      int                      `protobuf:"varint,1,req,name=desired,proto3" json:"desired" bson:"desired"`                                                //nolint:lll//nolint:lll
	Dealer       PlayerID                 `protobuf:"varint,2,req,name=playerID,proto3" json:"playerID" bson:"playerID"`                                             //nolint:lll//nolint:lll
	PlayerColors map[PlayerID]PlayerColor `protobuf:"map<varint,varint>,3,opt,name=playerColors,proto3" json:"playerColors,omitempty" bson:"playerColors,omitempty"` //nolint:lll//nolint:lll
}

type PlayerAction struct {
	GameID    GameID      `protobuf:"varint,1,req,name=gameID,proto3" json:"gameID" bson:"gameID"`          //nolint:lll
	ID        PlayerID    `protobuf:"varint,2,req,name=id,proto3" json:"id" bson:"id"`                      //nolint:lll
	Overcomes Blocker     `protobuf:"varint,3,req,name=overcomes,proto3" json:"overcomes" bson:"overcomes"` //nolint:lll
	Action    interface{} `protobuf:"-" json:"action" bson:"action"`                                        //nolint:lll
}

type DealAction struct {
	NumShuffles int `protobuf:"varint,1,req,name=numShuffles,proto3" json:"numShuffles" bson:"numShuffles"` //nolint:lll
}

type BuildCribAction struct {
	Cards []Card `protobuf:"Card,rep,1,req,name=cards,proto3" json:"cards" bson:"cards"` //nolint:lll
}

type CutDeckAction struct {
	Percentage float64 `protobuf:"varint,1,req,name=percentage,proto3" json:"percentage" bson:"percentage"` //nolint:lll
}

type PegAction struct {
	Card  Card `protobuf:"Card,1,req,name=card,proto3" json:"card" bson:"card"`    //nolint:lll
	SayGo bool `protobuf:"bool,2,req,name=sayGo,proto3" json:"sayGo" bson:"sayGo"` //nolint:lll
}

type CountHandAction struct {
	Pts int `protobuf:"varint,1,req,name=pts,proto3" json:"pts" bson:"pts"` //nolint:lll
}

type CountCribAction struct {
	Pts int `protobuf:"varint,1,req,name=pts,proto3" json:"pts" bson:"pts"` //nolint:lll
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
	ID              GameID                   `protobuf:"varint,1,name=id,proto3" json:"id" bson:"id"`                                                                   //nolint:lll
	Players         []Player                 `protobuf:"bytes,3,rep,name=players,proto3,proto3" json:"players" bson:"players"`                                          //nolint:lll
	Deck            Deck                     `protobuf:"-" json:"-" bson:"-"`                                                                                           //nolint:lll
	BlockingPlayers map[PlayerID]Blocker     `protobuf:"map<varint, varint>,4,opt,name=blockingPlayers,proto3" json:"blockingPlayers,omitempty" bson:"blockingPlayers"` //nolint:lll
	CurrentDealer   PlayerID                 `protobuf:"varint,5,opt,name=currentDealer,proto3" json:"currentDealer" bson:"currentDealer"`                              //nolint:lll
	PlayerColors    map[PlayerID]PlayerColor `protobuf:"map<varint, varint>,6,name=playerColors,proto3" json:"playerColors,omitempty" bson:"playerColors"`              //nolint:lll
	CurrentScores   map[PlayerColor]int      `protobuf:"map<varint, varint>,7,opt,name=currentScores,proto3" json:"currentScores" bson:"currentScores"`                 //nolint:lll
	LagScores       map[PlayerColor]int      `protobuf:"map<varint, varint>,8,opt,name=lagScores,proto3" json:"lagScores" bson:"lagScores"`                             //nolint:lll
	Phase           Phase                    `protobuf:"varint,9,opt,name=phase,proto3" json:"phase" bson:"phase"`                                                      //nolint:lll
	Hands           map[PlayerID][]Card      `protobuf:"map<varint, bytes>,10,opt,name=hands,proto3" json:"hands,omitempty" bson:"hands"`                               //nolint:lll
	CutCard         Card                     `protobuf:"Card,11,opt,name=cutCard,proto3" json:"cutCard,omitempty" bson:"cutCard"`                                       //nolint:lll
	Crib            []Card                   `protobuf:"Card,12,rep,name=crib,proto3" json:"crib,omitempty" bson:"crib"`                                                //nolint:lll
	PeggedCards     []PeggedCard             `protobuf:"PeggedCard,13,rep,opt,name=peggedCards,proto3" json:"peggedCards,omitempty" bson:"peggedCards"`                 //nolint:lll

	actions []PlayerAction
}
