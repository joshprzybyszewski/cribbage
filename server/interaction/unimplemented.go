package interaction

import "github.com/joshprzybyszewski/cribbage/model"

var _ Player = unimplemented{}

type unimplemented struct {
	myID model.PlayerID
}

func newUnimplemented(id model.PlayerID) Player {
	return unimplemented{
		myID: id,
	}
}

func (u unimplemented) ID() model.PlayerID {
	return u.myID
}

func (u unimplemented) NotifyBlocking(model.Blocker, model.Game, string) error {
	return nil
}

func (u unimplemented) NotifyMessage(model.Game, string) error {
	return nil
}

func (u unimplemented) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}
