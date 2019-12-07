// +build !prod

package interaction

import (
	"fmt"

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
	if e == nil {
		return nil
	}
	fmt.Printf("NotifyBlocking | b, g, s := %+v, %+v, %+v\n", b, g, s)
	return nil
}
func (e *Empty) NotifyMessage(g model.Game, s string) error {
	if e == nil {
		return nil
	}
	fmt.Printf("NotifyMessage | g, s := %+v, %+v\n", g, s)
	return nil
}
func (e *Empty) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	if e == nil {
		return nil
	}
	fmt.Printf("NotifyScoreUpdate | g, msgs := %+v, %+v\n", g, msgs)
	return nil
}
