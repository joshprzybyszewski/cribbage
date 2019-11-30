package interaction

import "github.com/joshprzybyszewski/cribbage/model"

var _ npcPlayer = (*dumbNPCPlayer)(nil)
var _ Player = (*dumbNPCPlayer)(nil)

type dumbNPCPlayer struct{}

func (npc *dumbNPCPlayer) ID() model.PlayerID {
	return `dumbNPC`
}

// Methods satisfying Player interface
func (npc *dumbNPCPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return handleNPCBlocker(npc, b, g, npc.ID())
}
func (npc *dumbNPCPlayer) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (npc *dumbNPCPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}

// Methods satisfying npcPlayer interface
func (npc *dumbNPCPlayer) buildCrib(g model.Game) model.BuildCribAction {
	return model.BuildCribAction{}
}
func (npc *dumbNPCPlayer) peg(g model.Game) model.PegAction {
	return model.PegAction{}
}
