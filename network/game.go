package network

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
)

func ConvertToGetGameResponse(g model.Game) GetGameResponse {
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

func ConvertToGetGameResponseForPlayer(g model.Game, pID model.PlayerID) (GetGameResponse, error) {
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
	resp := ConvertToGetGameResponse(g)
	resp.Hands = getRevealedCards(g, pID)
	if g.Phase >= model.CribCounting {
		resp.Crib = convertCards(g.Crib)
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
		rev[c.PlayerID] = append(rev[c.PlayerID], newCardFromModel(c.Card))
	}
	rev[me] = convertCards(g.Hands[me])
	return rev
}

func ConvertToCreateGameResponse(g model.Game) CreateGameResponse {
	return CreateGameResponse{
		ID:              g.ID,
		Players:         newPlayersFromModels(g.Players),
		PlayerColors:    convertColors(g.PlayerColors),
		BlockingPlayers: convertBlockingPlayers(g.BlockingPlayers),
		CurrentDealer:   g.CurrentDealer,
	}
}
