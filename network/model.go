package network

import "github.com/joshprzybyszewski/cribbage/model"

type CreateInteractionRequest struct {
	PlayerID      model.PlayerID `json:"playerID"`
	LocalhostPort string         `json:"localhost_port,omitempty"`
	NPCType       model.PlayerID `json:"npc_type,omitempty"`
}

type GetGameResponse struct {
	ID              model.GameID              `json:"id"`
	Players         []Player                  `json:"players"`
	PlayerColors    map[model.PlayerID]string `json:"player_colors,omitempty"`
	CurrentScores   map[string]int            `json:"current_scores"`
	LagScores       map[string]int            `json:"lag_scores"`
	Phase           string                    `json:"phase"`
	BlockingPlayers map[model.PlayerID]string `json:"blocking_players,omitempty"`
	CurrentDealer   model.PlayerID            `json:"current_dealer"`

	// Scrub out the other players' hands when returning this
	Hands map[model.PlayerID][]Card `json:"hands,omitempty"`
	// Scrub out the crib until the phase is correct
	Crib []Card `json:"crib,omitempty"`

	CutCard     Card   `json:"cut_card"`
	PeggedCards []Card `json:"pegged_cards,omitempty"`
}

type CreateGameRequest struct {
	PlayerIDs []model.PlayerID `json:"playerIDs"`
}

type CreateGameResponse struct {
	ID              model.GameID              `json:"id"`
	Players         []Player                  `json:"players"`
	PlayerColors    map[model.PlayerID]string `json:"player_colors,omitempty"`
	BlockingPlayers map[model.PlayerID]string `json:"blocking_players,omitempty"`
	CurrentDealer   model.PlayerID            `json:"current_dealer"`
}

type Player struct {
	ID   model.PlayerID `json:"id"`
	Name string         `json:"name"`
}

type CreatePlayerRequest struct {
	Player Player `json:"player"`
}

type GetPlayerResponse struct {
	Player Player                             `json:"player"`
	Games  map[model.GameID]model.PlayerColor `json:"games"`
}

type CreatePlayerResponse struct {
	Player Player `json:"player"`
}

type Card struct {
	Suit  string
	Value int
	Name  string
}
