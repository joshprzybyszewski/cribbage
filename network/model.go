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
type GetGameResponse struct {
	ID              model.GameID                         `json:"id"`
	Players         []Player                             `json:"players"`
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
}

func NewGetGameResponse(g model.Game) GetGameResponse {
	return GetGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    g.PlayerColors,
		CurrentScores:   g.CurrentScores,
		LagScores:       g.LagScores,
		Phase:           g.Phase,
		BlockingPlayers: g.BlockingPlayers,
		CurrentDealer:   g.CurrentDealer,
		Hands:           g.Hands,
		Crib:            g.Crib,
		CutCard:         g.CutCard,
		PeggedCards:     g.PeggedCards,
	}
}

// TODO figure out the minimum info the client will need
type CreateGameResponse struct {
	ID              model.GameID                         `json:"id"`
	Players         []Player                             `json:"players"`
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
}

func NewCreateGameResponse(g model.Game) CreateGameResponse {
	return CreateGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    g.PlayerColors,
		CurrentScores:   g.CurrentScores,
		LagScores:       g.LagScores,
		Phase:           g.Phase,
		BlockingPlayers: g.BlockingPlayers,
		CurrentDealer:   g.CurrentDealer,
		Hands:           g.Hands,
		Crib:            g.Crib,
		CutCard:         g.CutCard,
		PeggedCards:     g.PeggedCards,
	}
}

type Player struct {
	ID   model.PlayerID `json:"id"`
	Name string         `json:"name"`
}

func newPlayerFromModel(p model.Player) Player {
	return Player{
		ID:   p.ID,
		Name: p.Name,
	}
}

func newPlayersFromModels(pms []model.Player) []Player {
	ps := make([]Player, len(pms))
	for i, pm := range pms {
		ps[i] = newPlayerFromModel(pm)
	}
	return ps
}

type CreatePlayerRequest struct {
	Player Player
}

type GetPlayerResponse struct {
	Player Player
	Games  map[model.GameID]model.PlayerColor `json:"games"`
}

func NewGetPlayerResponseFromModel(pm model.Player) GetPlayerResponse {
	return GetPlayerResponse{
		Player: Player{
			ID:   pm.ID,
			Name: pm.Name,
		},
		Games: pm.Games,
	}
}

type CreatePlayerResponse struct {
	Player Player
}

func NewCreatePlayerResponseFromModel(pm model.Player) CreatePlayerResponse {
	return CreatePlayerResponse{
		Player: Player{
			ID:   pm.ID,
			Name: pm.Name,
		},
	}
}
