package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ Player = (*empty)(nil)

type empty struct {
	PID model.PlayerID
}

func Empty(pID model.PlayerID) Player {
	return &empty{
		PID: pID,
	}
}

func (e *empty) ID() model.PlayerID {
	if e == nil {
		return model.InvalidPlayerID
	}
	return e.PID
}
func (e *empty) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return nil
}
func (e *empty) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (e *empty) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}
