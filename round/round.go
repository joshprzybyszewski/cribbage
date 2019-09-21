package round

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
)

type Round struct {
	deck         *cards.Deck
	players      []*Player
	dealerIndex  int
	CurrentStage Stage

	// The map of player to list of cards they've placed in the crib
	// for three player games, use nil for the dealt card
	cribCards map[*Player][]cards.Card

	// the ordered list of cards which is added to as players play a card during pegging
	peggedCards []cards.Card

	// the number we are currently at in pegging
	currentPeg int
}

func NewGame(players []*Player, dealerIndex int) *Round {
	if len(players) > 4 || len(players) < 2 {
		// cannot play cribbage with any number besides 2, 3, and 4
		return nil
	}

	return &Round{
		deck:         cards.NewDeck(),
		players:      players,
		dealerIndex:  dealerIndex,
		CurrentStage: Deal,
		cribCards:    map[int][]cards.Card{},
		peggedCards:  make([]cards.Card, 0, 4*len(players)),
		currentPeg:   0,
	}
}

func NextRound(prev *Round) *Round {
	if prev == nil {
		return nil
	}

	nextDealerIndex := (prev.dealerIndex + 1) % len(prev.players)
	prev.deck.Shuffle()

	return &Round{
		deck:         prev.deck,
		players:      prev.players,
		dealerIndex:  nextDealerIndex,
		CurrentStage: Deal,
		cribCards:    map[int][]cards.Card{},
		peggedCards:  prev.peggedCards[:0],
		currentPeg:   0,
	}
}

func (r *Round) ShuffleDeck() {
	if r.CurrentStage != Deal {
		return
	}

	r.deck.Shuffle()
}

func (r *Round) DealCards() error {
	if r.CurrentStage != Deal {
		return errors.New(`cannot deal in current stage`)
	}

	switch len(r.players) {
	case 2:
		return r.dealTwoPlayerGame()
	case 3:
		return r.dealThreePlayerGame()
	case 4:
		return r.dealFourPlayerGame()
	}

	return errors.New(`should not be in this game`)
}

func (r *Round) dealTwoPlayerGame() error {
	// need to deal 12 cards; 6 to each player
	for i := 0; i < 12; i++ {
		// this ensures we always deal to the other player first, alternating on each deal
		// so if the dealer index is 0,  then it should go 1,0,1,0,1,0,1,0,1,0,1,0
		// and if the dealer index is 1, then it should go 0,1,0,1,0,1,0,1,0,1,0,1
		pi := 1 - (r.dealerIndex - (i % 2))
		err := r.players[pi].AcceptCard(r.deck.Deal())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Round) dealThreePlayerGame() error {
	return errors.New(`need to implement 3 player game`)
}

func (r *Round) dealFourPlayerGame() error {
	return errors.New(`need to implement 4 player game`)
}

func (r *Round) AddToCrib(p *player, c cards.Card) error {
	if r.CurrentStage != BuildCrib {
		return errors.New(`cannot add to crib in current stage`)
	}

	r.cribCards[p] = append(r.cribCards[p], c)
	return nil
}

func (r *Round) Cut(perc float64) error {
	if r.CurrentStage != Cut {
		return errors.New(`cannot cut in current stage`)
	}

	// TODO cut

	return errors.New(`unimplemented`)

}

func (r *Round) PegCard(player *Player, c cards.Card) (pts int, err error) {
	return 0, errors.New(`unimplemented`)
}

func (r *Round) SayGo(player *Player) (err error) {
	return errors.New(`need to validate and stuff`)
}
