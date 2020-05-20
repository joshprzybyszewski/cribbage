package network

import "github.com/joshprzybyszewski/cribbage/model"

const EmptyInteraction = ``

type CreateInteractionRequest struct {
	PlayerID      model.PlayerID `json:"playerID"`
	LocalhostPort string         `json:"localhost_port,omitempty"`
	NPCType       string         `json:"npc_type,omitempty"`
}
