package mysql

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.GameService = (*gameService)(nil)

type gameService struct{}

func (gs *gameService) Get(id model.GameID) (model.Game, error) {
	return model.Game{}, nil
}

func (gs *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	return model.Game{}, nil
}

func (gs *gameService) UpdatePlayerColor(gID model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	return nil
}

func (gs *gameService) Save(g model.Game) error {
	return nil
}

func decodeHand(cards uint64, idx int) model.Card {
	// we need 6 bits per card to store 0-51
	bitsToShift := uint(idx * 6)
	cardRank := int((cards >> bitsToShift) & 0x3F)
	return model.NewCardFromNumber(cardRank)
}
