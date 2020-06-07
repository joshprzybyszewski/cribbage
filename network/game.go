package network

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

func NewGetGameResponse(g model.Game) GetGameResponse {
	currentScores, lagScores := convertScores(g.CurrentScores, g.LagScores)
	return GetGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    convertColors(g.PlayerColors),
		CurrentScores:   currentScores,
		LagScores:       lagScores,
		Phase:           g.Phase.String(),
		BlockingPlayers: convertBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
		CutCard:         newCardFromModel(g.CutCard),
		PeggedCards:     convertPeggedCards(g.PeggedCards),
	}
}

func NewGetGameResponseForPlayer(g model.Game, pID model.PlayerID) (GetGameResponse, error) {
	pIsInGame := false
	for _, p := range g.Players {
		if p.ID == pID {
			pIsInGame = true
			break
		}
	}
	if !pIsInGame {
		return GetGameResponse{}, errors.New(`player does not exist in game`)
	}
	resp := NewGetGameResponse(g)
	resp.Hands = convertHands(g.Hands)
	if g.Phase < model.Counting {
		resp.Hands = map[model.PlayerID][]Card{
			pID: resp.Hands[pID],
		}
	}
	if g.Phase >= model.CribCounting {
		resp.Crib = convertCards(g.Crib)
	}
	return resp, nil
}

func NewCreateGameResponse(g model.Game) CreateGameResponse {
	return CreateGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    convertColors(g.PlayerColors),
		BlockingPlayers: convertBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
	}
}
