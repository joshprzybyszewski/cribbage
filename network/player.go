package network

import "github.com/joshprzybyszewski/cribbage/model"

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

type GetActiveGamesForPlayerResponse struct {
	Player      Player                  `json:"player"`
	ActiveGames map[model.GameID]string `json:"activeGames"`
}

func ConvertToGetActiveGamesForPlayerResponse(p model.Player, games map[model.GameID]model.Game) GetActiveGamesForPlayerResponse {
	return GetActiveGamesForPlayerResponse{
		Player: Player{
			ID:   p.ID,
			Name: p.Name,
		},
		ActiveGames: convertToParticipatingGames(p, games),
	}
}

func convertToParticipatingGames(p model.Player, games map[model.GameID]model.Game) map[model.GameID]string {
	res := make(map[model.GameID]string, len(p.Games))
	for gID := range p.Games {
		if mg, ok := games[gID]; ok {
			res[gID] = getPlayerNames(mg)
		}
	}
	return res
}

func getPlayerNames(mg model.Game) string {
	res := ``
	for i, p := range mg.Players {
		if i > 0 {
			res += `, `
		}
		res += p.Name
	}
	return res
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
