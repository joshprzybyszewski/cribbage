package network

import "github.com/joshprzybyszewski/cribbage/model"

type CreateInteractionRequest struct {
	PlayerID model.PlayerID `json:"playerID"`
	Mode     string         `json:"mode"`
	Info     interface{}    `json:"info"`
}
