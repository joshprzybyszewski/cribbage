package network

import "github.com/joshprzybyszewski/cribbage/model"

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

func NewGetPlayerResponseFromModel(pm model.Player) GetPlayerResponse {
	return GetPlayerResponse{
		Player: Player{
			ID:   pm.ID,
			Name: pm.Name,
		},
		Games: convertGames(pm.Games),
	}
}

func convertGames(modelGames map[model.GameID]model.PlayerColor) map[model.GameID]string {
	games := make(map[model.GameID]string, len(modelGames))
	for g, c := range modelGames {
		games[g] = c.String()
	}
	return games
}

func NewCreatePlayerResponseFromModel(pm model.Player) CreatePlayerResponse {
	return CreatePlayerResponse{
		Player: Player{
			ID:   pm.ID,
			Name: pm.Name,
		},
	}
}
