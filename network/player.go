package network

import (
	"time"

	"github.com/joshprzybyszewski/cribbage/model"
)

type Player struct {
	ID   model.PlayerID `json:"id"`
	Name string         `json:"name"`
}

type CreatePlayerRequest struct {
	Player Player `json:"player"`
}

type CreatePlayerResponse struct {
	Player Player `json:"player"`
}

func ConvertToCreatePlayerResponse(pm model.Player) CreatePlayerResponse {
	return CreatePlayerResponse{
		Player: Player{
			ID:   pm.ID,
			Name: pm.Name,
		},
	}
}

type GetPlayerResponse struct {
	Player Player `json:"player"`
}

func ConvertToGetPlayerResponse(p model.Player) GetPlayerResponse {
	return GetPlayerResponse{
		Player: Player{
			ID:   p.ID,
			Name: p.Name,
		},
	}
}

type ActiveGame struct {
	PlayerNamesByID  map[model.PlayerID]string `json:"players"`
	PlayerColorsByID map[model.PlayerID]string `json:"colors"`
	Created          time.Time                 `json:"created"`
	LastMove         time.Time                 `json:"lastMove"`
}

type GetActiveGamesForPlayerResponse struct {
	Player      Player                      `json:"player"`
	ActiveGames map[model.GameID]ActiveGame `json:"activeGames"`
}

func ConvertToGetActiveGamesForPlayerResponse(
	p model.Player,
	games map[model.GameID]model.Game,
) GetActiveGamesForPlayerResponse {

	return GetActiveGamesForPlayerResponse{
		Player: Player{
			ID:   p.ID,
			Name: p.Name,
		},
		ActiveGames: convertToParticipatingGames(p, games),
	}
}

func convertToParticipatingGames(p model.Player, games map[model.GameID]model.Game) map[model.GameID]ActiveGame {
	res := make(map[model.GameID]ActiveGame, len(p.Games))
	for gID := range p.Games {
		if mg, ok := games[gID]; ok {
			res[gID] = getActiveGame(mg)
		}
	}
	return res
}

func getActiveGame(mg model.Game) ActiveGame {
	ag := ActiveGame{
		PlayerNamesByID:  map[model.PlayerID]string{},
		PlayerColorsByID: map[model.PlayerID]string{},
	}
	for _, p := range mg.Players {
		pID := p.ID
		ag.PlayerNamesByID[pID] = p.Name
		ag.PlayerColorsByID[pID] = mg.PlayerColors[pID].String()
	}
	return ag
}

func convertToPlayers(pms []model.Player) []Player {
	ps := make([]Player, len(pms))
	for i, pm := range pms {
		ps[i] = convertToPlayer(pm)
	}
	return ps
}

func convertFromPlayers(pms []Player) []model.Player {
	ps := make([]model.Player, len(pms))
	for i, pm := range pms {
		ps[i] = convertFromPlayer(pm)
	}
	return ps
}

func convertToPlayer(p model.Player) Player {
	return Player{
		ID:   p.ID,
		Name: p.Name,
	}
}

func convertFromPlayer(p Player) model.Player {
	return model.Player{
		ID:   p.ID,
		Name: p.Name,
	}
}
