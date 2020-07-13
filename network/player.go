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

type ActiveGamePlayer struct {
	ID    model.PlayerID `json:"id"`
	Name  string         `json:"name"`
	Color string         `json:"color"`
}

type ActiveGame struct {
	GameID   model.GameID       `json:"gameID"`
	Players  []ActiveGamePlayer `json:"players"`
	Created  time.Time          `json:"created"`
	LastMove time.Time          `json:"lastMove"`
}

type GetActiveGamesForPlayerResponse struct {
	Player      Player       `json:"player"`
	ActiveGames []ActiveGame `json:"activeGames"`
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

func convertToParticipatingGames(
	p model.Player,
	games map[model.GameID]model.Game,
) []ActiveGame {
	res := make([]ActiveGame, 0, len(p.Games))
	for gID := range p.Games {
		if mg, ok := games[gID]; ok {
			res = append(res, getActiveGame(mg))
		}
	}
	return res
}

func getActiveGame(mg model.Game) ActiveGame {
	ag := ActiveGame{
		GameID: mg.ID,
	}
	ag.Players = make([]ActiveGamePlayer, len(mg.Players))
	for i, p := range mg.Players {
		pID := p.ID
		ag.Players[i] = ActiveGamePlayer{
			ID:    pID,
			Name:  p.Name,
			Color: mg.PlayerColors[pID].String(),
		}
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

func convertTeamsToPlayersAndPlayerColors(
	ts []GetGameResponseTeam,
) ([]model.Player, map[model.PlayerID]model.PlayerColor) {
	// length of 4 at most
	players := make([]model.Player, 0)
	playerColors := make(map[model.PlayerID]model.PlayerColor, len(ts))
	for _, t := range ts {
		for _, p := range t.Players {
			players = append(players, convertFromPlayer(p))
			playerColors[p.ID] = model.NewPlayerColorFromString(t.Color)
		}
	}
	return players, playerColors
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
