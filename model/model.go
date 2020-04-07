package model

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit Suit `protobuf:"varint,1,req,name=suit,proto3" json:"s" bson:"s"` //nolint:lll
	// Ace is 1, King is 13
	Value int `protobuf:"varint,2,req,name=value,proto3" json:"v" bson:"v"` //nolint:lll
}

const NumCardsPerDeck = 52
const JackValue = 11 // Ace is 1, King is 13

type PeggedCard struct {
	Card `protobuf:"-" json:"pc" bson:"pc"` //nolint:lll

	Action   int      `protobuf:"-" json:"aIdx" bson:"aIdx"` //nolint:lll
	PlayerID PlayerID `protobuf:"-" json:"pID" bson:"pID"`   //nolint:lll
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
	Auth
	ID    PlayerID               `protobuf:"varint,1,req,name=id,proto3" json:"id" bson:"id"`                           //nolint:lll
	Name  string                 `protobuf:"string,2,req,name=name,proto3" json:"n" bson:"n"`                           //nolint:lll
	Games map[GameID]PlayerColor `protobuf:"map<varint, varint>,3,opt,name=games,proto3" json:"gs,omitempty" bson:"gs"` //nolint:lll
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
	Desired      int                      `protobuf:"varint,1,req,name=desired,proto3" json:"d" bson:"d"`                                        //nolint:lll//nolint:lll
	Dealer       PlayerID                 `protobuf:"varint,2,req,name=playerID,proto3" json:"pID" bson:"pID"`                                   //nolint:lll//nolint:lll
	PlayerColors map[PlayerID]PlayerColor `protobuf:"map<varint,varint>,3,opt,name=playerColors,proto3" json:"pc,omitempty" bson:"pc,omitempty"` //nolint:lll//nolint:lll
}

type PlayerAction struct {
	GameID    GameID      `protobuf:"varint,1,req,name=gameID,proto3" json:"gID" bson:"gID"` //nolint:lll
	ID        PlayerID    `protobuf:"varint,2,req,name=id,proto3" json:"pID" bson:"pID"`     //nolint:lll
	Overcomes Blocker     `protobuf:"varint,3,req,name=overcomes,proto3" json:"o" bson:"o"`  //nolint:lll
	Action    interface{} `protobuf:"-" json:"a" bson:"a"`                                   //nolint:lll
}

type DealAction struct {
	NumShuffles int `protobuf:"varint,1,req,name=numShuffles,proto3" json:"ns" bson:"ns"` //nolint:lll
}

type BuildCribAction struct {
	Cards []Card `protobuf:"Card,rep,1,req,name=cards,proto3" json:"cs" bson:"cs"` //nolint:lll
}

type CutDeckAction struct {
	Percentage float64 `protobuf:"varint,1,req,name=percentage,proto3" json:"p" bson:"p"` //nolint:lll
}

type PegAction struct {
	Card  Card `protobuf:"Card,1,req,name=card,proto3" json:"c" bson:"c"`    //nolint:lll
	SayGo bool `protobuf:"bool,2,req,name=sayGo,proto3" json:"sg" bson:"sg"` //nolint:lll
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

func (p Phase) String() string { //nolint:gocyclo
	switch p {
	case Deal:
		return `Deal`
	case BuildCribReady:
		return `BuildCribReady`
	case BuildCrib:
		return `BuildCrib`
	case CutReady:
		return `CutReady`
	case Cut:
		return `Cut`
	case PeggingReady:
		return `PeggingReady`
	case Pegging:
		return `Pegging`
	case CountingReady:
		return `CountingReady`
	case Counting:
		return `Counting`
	case CribCountingReady:
		return `CribCountingReady`
	case CribCounting:
		return `CribCounting`
	case DealingReady:
		return `DealingReady`
	}
	return `unknown`
}

const (
	WinningScore    int = 121
	MaxPeggingValue int = 31
)

const (
	MinPlayerGame int = 2
	MaxPlayerGame int = 4
)

// Game represents all of the data needed for a game of cribbage
// between 2, 3, or 4 players
type Game struct {
	// The unique identifier used to reference this game
	ID GameID `protobuf:"-" json:"id" bson:"id"` //nolint:lll

	// The players playing this game and their colors
	Players      []Player                 `protobuf:"-" json:"ps" bson:"ps"`             //nolint:lll
	PlayerColors map[PlayerID]PlayerColor `protobuf:"-" json:"pcs,omitempty" bson:"pcs"` //nolint:lll

	// The current (and lagging) scores
	CurrentScores map[PlayerColor]int `protobuf:"-" json:"cs" bson:"cs"` //nolint:lll
	LagScores     map[PlayerColor]int `protobuf:"-" json:"ls" bson:"ls"` //nolint:lll

	// What phase this game is in
	Phase Phase `protobuf:"-" json:"p" bson:"p"` //nolint:lll
	// Who is blocking and why
	BlockingPlayers map[PlayerID]Blocker `protobuf:"-" json:"bps,omitempty" bson:"bps"` //nolint:lll

	// The identifier for the current dealer
	CurrentDealer PlayerID `protobuf:"-" json:"cd" bson:"cd"` //nolint:lll

	// The hands of each player
	Hands map[PlayerID][]Card `protobuf:"-" json:"hs,omitempty" bson:"hs"` //nolint:lll
	// The cards currently in the crib
	Crib []Card `protobuf:"-" json:"c,omitempty" bson:"c"` //nolint:lll

	// The flipped card which acts as the lead
	CutCard Card `protobuf:"-" json:"cc" bson:"cc"` //nolint:lll

	// An ordered list of previously pegged cards (which includes who pegged them), most recent last
	PeggedCards []PeggedCard `protobuf:"-" json:"pegged,omitempty" bson:"pegged"` //nolint:lll

	// An ordered list of player actions
	Actions []PlayerAction `protobuf:"-" json:"as" bson:"as"` //nolint:lll
}
