package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

// NPC is an enum specifying which type of NPC
type NPC int

// Dumb, Simple, and Calculated are supported
const (
	Dumb NPC = iota
	Simple
	Calculated
)

var npcIDs = [...]string{
	Dumb:       `dumbNPC`,
	Simple:     `simpleNPC`,
	Calculated: `calculatedNPC`,
}

var _ Player = (*npcPlayer)(nil)

type npcPlayer struct {
	Type NPC
}

func (npc *npcPlayer) ID() model.PlayerID {
	return model.PlayerID(npcIDs[npc.Type])
}

func (npc *npcPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return handleNPCBlocker(npc.Type, b, g)
}
func (npc *npcPlayer) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (npc *npcPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}
