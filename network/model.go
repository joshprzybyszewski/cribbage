package network

import "github.com/joshprzybyszewski/cribbage/model"

type CreatePlayerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type CreateGameRequest struct {
	PlayerIDs []model.PlayerID `json:"playerIDs"`
}

type CreateInteractionRequest struct {
	PlayerID model.PlayerID `json:"playerID"`
	Mode     string         `json:"mode"`
	Info     interface{}    `json:"info"`
}
