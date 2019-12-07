// +build !prod

package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ Player = (*Empty)(nil)

type Empty struct {
	PID model.PlayerID
}

func (e *Empty) ID() model.PlayerID {
	if e == nil {
		return model.InvalidPlayerID
	}
	return e.PID
}
func (e *Empty) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return nil
}
func (e *Empty) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (e *Empty) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}
