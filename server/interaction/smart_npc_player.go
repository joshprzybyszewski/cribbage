package interaction

import "github.com/joshprzybyszewski/cribbage/model"

var _ Player = (*smartNPCPlayer)(nil)

type smartNPCPlayer struct{}

func (npc *smartNPCPlayer) ID() model.PlayerID {
	return `smartNPC`
}

func (npc *smartNPCPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return nil
}

func (npc *smartNPCPlayer) NotifyMessage(g model.Game, msg string) error {
	return nil
}

func (npc *smartNPCPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}
