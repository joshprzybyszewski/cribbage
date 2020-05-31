package network

import "github.com/joshprzybyszewski/cribbage/model"

type CreateInteractionRequest struct {
	PlayerID      model.PlayerID `json:"playerID"`
	LocalhostPort string         `json:"localhost_port,omitempty"`
	NPCType       model.PlayerID `json:"npc_type,omitempty"`
}

type CreateGameRequest struct {
	PlayerIDs []model.PlayerID `json:"playerIDs"`
}

// TODO figure out the minimum info the client will need
type CreateGameResponse struct {
	Players         []model.Player                       `json:"players"`
	PlayerColors    map[model.PlayerID]model.PlayerColor `json:"player_colors,omitempty"`
	CurrentScores   map[model.PlayerColor]int            `json:"current_scores"`
	LagScores       map[model.PlayerColor]int            `json:"lag_scores"`
	Phase           model.Phase                          `json:"phase"`
	BlockingPlayers map[model.PlayerID]model.Blocker     `json:"blocking_players,omitempty"`
	CurrentDealer   model.PlayerID                       `json:"current_dealer"`
	Hands           map[model.PlayerID][]model.Card      `json:"hands,omitempty"`
	Crib            []model.Card                         `json:"crib,omitempty"`
	CutCard         model.Card                           `json:"cut_card"`
	PeggedCards     []model.PeggedCard                   `json:"pegged_cards,omitempty"`
	Actions         []model.PlayerAction                 `json:"actions"`
}

type CreatePlayerRequest struct {
	ID   model.PlayerID `json:"id"`
	Name string         `json:"name"`
}

type CreatePlayerResponse struct {
	ID   model.PlayerID `json:"id"`
	Name string         `json:"name"`
}
