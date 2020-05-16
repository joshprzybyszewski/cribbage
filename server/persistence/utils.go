package persistence

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

func ValidateLatestActionBelongs(mg model.Game) error {
	if mg.NumActions() == 0 {
		return nil
	}

	la := mg.Actions[mg.NumActions()-1]
	if la.GameID != mg.ID {
		return ErrGameActionWrongGame
	}
	found := false
	for _, p := range mg.Players {
		if p.ID == la.ID {
			found = true
			break
		}
	}
	if !found {
		return ErrGameActionWrongPlayer
	}

	return nil
}
