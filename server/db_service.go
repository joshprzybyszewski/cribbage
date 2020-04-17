package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

type DBService struct {
	db persistence.DB
}

func NewDBService(db persistence.DB) *DBService {
	return &DBService{
		db: db,
	}
}

func (service *DBService) CreatePlayer(id model.PlayerID, displayName string) (model.Player, error) {
	player := model.Player{
		ID:    model.PlayerID(id),
		Name:  displayName,
		Games: make(map[model.GameID]model.PlayerColor),
	}
	err := service.db.CreatePlayer(player)
	if err != nil {
		return model.Player{}, err
	}
	return player, nil
}

func (service *DBService) CreateGame(pIDs []model.PlayerID) (model.Game, error) {
	players := make([]model.Player, len(pIDs))
	for i, id := range pIDs {
		p, err := service.db.GetPlayer(id)
		if err != nil {
			return model.Game{}, err
		}
		players[i] = p
	}

	pAPIs, err := service.getPlayerAPIs(players)
	if err != nil {
		return model.Game{}, err
	}

	mg, err := play.CreateGame(players, pAPIs)
	if err != nil {
		return model.Game{}, err
	}

	err = service.db.SaveGame(mg)
	if err != nil {
		return model.Game{}, err
	}

	for _, pID := range pIDs {
		err = service.db.AddPlayerColorToGame(pID, mg.PlayerColors[pID], mg.ID)
		if err != nil {
			return model.Game{}, err
		}
	}

	return mg, nil
}

func (service *DBService) getPlayerAPIs(players []model.Player) (map[model.PlayerID]interaction.Player, error) {
	pAPIs := make(map[model.PlayerID]interaction.Player, len(players))
	for _, p := range players {
		var pAPI interaction.Player
		pm, err := service.db.GetInteraction(p.ID)

		for i, m := range pm.Interactions {
			if m.Mode == interaction.NPC {
				m.Info = &npcActionHandler{}
				pm.Interactions[i] = m
			}
		}

		if err != nil {
			if err != persistence.ErrInteractionNotFound {
				return nil, err
			}
			pAPI = interaction.Empty(p.ID)
		} else {
			pAPI, err = interaction.FromPlayerMeans(pm)
			if err != nil {
				return nil, err
			}
		}

		pAPIs[p.ID] = pAPI
	}
	return pAPIs, nil
}
