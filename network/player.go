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
	Player Player                  `json:"player"`
	Games  map[model.GameID]string `json:"games"`
}

func ConvertToGetPlayerResponse(p model.Player) GetPlayerResponse {
	return GetPlayerResponse{
		Player: Player{
			ID:   p.ID,
			Name: p.Name,
		},
		Games: convertToParticipatingGames(p.Games),
	}
}

func convertToParticipatingGames(mgs map[model.GameID]model.PlayerColor) map[model.GameID]string {
	games := make(map[model.GameID]string, len(mgs))
	for g, c := range mgs {
		games[g] = c.String()
	}
	return games
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
