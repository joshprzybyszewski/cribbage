package network

import "github.com/joshprzybyszewski/cribbage/model"

type CreatePlayerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type CreateGameRequest struct {
	PlayerIDs []model.PlayerID `json:"playerIDs"`
}
