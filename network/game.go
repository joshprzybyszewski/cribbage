package network

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

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

func ConvertToCreateGameResponse(g model.Game) CreateGameResponse {
	return CreateGameResponse{
		ID:              g.ID,
		Players:         convertToPlayers(g.Players),
		PlayerColors:    convertToColors(g.PlayerColors),
		BlockingPlayers: convertToBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
	}
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
	Hands           map[model.PlayerID][]Card `json:"hands,omitempty"`
	Crib            []Card                    `json:"crib,omitempty"`
	CutCard         Card                      `json:"cut_card"`
	PeggedCards     []PeggedCard              `json:"pegged_cards,omitempty"`
}

func ConvertToGetGameResponse(g model.Game) GetGameResponse {
	currentScores, lagScores := convertToScores(g.CurrentScores, g.LagScores)
	return GetGameResponse{
		ID:              g.ID,
		Players:         convertToPlayers(g.Players),
		PlayerColors:    convertToColors(g.PlayerColors),
		CurrentScores:   currentScores,
		LagScores:       lagScores,
		Phase:           g.Phase.String(),
		BlockingPlayers: convertToBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
		CutCard:         convertToCard(g.CutCard),
		PeggedCards:     convertPeggedCards(g.PeggedCards),
	}
}

func ConvertToGetGameResponseForPlayer(g model.Game, pID model.PlayerID) (GetGameResponse, error) {
	playing := false
	for _, p := range g.Players {
		if p.ID == pID {
			playing = true
			break
		}
	}
	if !playing {
		return GetGameResponse{}, errors.New(`player does not exist in game`)
	}
	resp := ConvertToGetGameResponse(g)
	resp.Hands = getRevealedCards(g, pID)
	if g.Phase >= model.CribCounting {
		resp.Crib = convertToCards(g.Crib)
	}
	return resp, nil
}

func getRevealedCards(g model.Game, me model.PlayerID) map[model.PlayerID][]Card {
	rev := make(map[model.PlayerID][]Card, len(g.Players))
	for pID := range rev {
		// we don't know how many cards will be revealed, but we know it's a max of 4
		rev[pID] = make([]Card, 0, 4)
	}
	for _, c := range g.PeggedCards {
		if c.PlayerID == me {
			continue
		}
		rev[c.PlayerID] = append(rev[c.PlayerID], convertToCard(c.Card))
	}
	rev[me] = convertToCards(g.Hands[me])
	return rev
}
