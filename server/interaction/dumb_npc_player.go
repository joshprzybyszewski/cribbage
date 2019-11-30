package interaction

import "github.com/joshprzybyszewski/cribbage/model"

var _ Player = (*dumbNPCPlayer)(nil)

type dumbNPCPlayer struct{}

func (npc *dumbNPCPlayer) ID() model.PlayerID {
	return `dumbNPC`
}

func (npc *dumbNPCPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return nil
}

func (npc *dumbNPCPlayer) NotifyMessage(g model.Game, msg string) error {
	return nil
}

func (npc *dumbNPCPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}
