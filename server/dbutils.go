package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/play"
)

func createPlayer(db persistence.DB, id model.PlayerID, displayName string) (model.Player, error) {
	player := model.Player{
		ID:    model.PlayerID(id),
		Name:  displayName,
		Games: make(map[model.GameID]model.PlayerColor),
	}
	err := db.CreatePlayer(player)
	if err != nil {
		return model.Player{}, err
	}
	return player, nil
}

func createGame(db persistence.DB, pIDs []model.PlayerID) (model.Game, error) {
	players := make([]model.Player, len(pIDs))
	for i, id := range pIDs {
		p, err := db.GetPlayer(id)
		if err != nil {
			return model.Game{}, err
		}
		players[i] = p
	}

	pAPIs, err := getPlayerAPIs(db, players)
	if err != nil {
		return model.Game{}, err
	}

	mg, err := play.CreateGame(players, pAPIs)
	if err != nil {
		return model.Game{}, err
	}

	err = db.SaveGame(mg)
	if err != nil {
		return model.Game{}, err
	}

	for _, pID := range pIDs {
		err = db.AddPlayerColorToGame(pID, mg.PlayerColors[pID], mg.ID)
		if err != nil {
			return model.Game{}, err
		}
	}

	return mg, nil
}

// func (service *dbService) getPlayerAPIs(players []model.Player) (map[model.PlayerID]interaction.Player, error) {
// 	pAPIs := make(map[model.PlayerID]interaction.Player, len(players))
// 	for _, p := range players {
// 		var pAPI interaction.Player
// 		pm, err := service.db.GetInteraction(p.ID)

// 		for i, m := range pm.Interactions {
// 			if m.Mode == interaction.NPC {
// 				m.Info = &npcActionHandler{}
// 				pm.Interactions[i] = m
// 			}
// 		}

// 		if err != nil {
// 			if err != persistence.ErrInteractionNotFound {
// 				return nil, err
// 			}
// 			pAPI = interaction.Empty(p.ID)
// 		} else {
// 			pAPI, err = interaction.FromPlayerMeans(pm)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}

// 		pAPIs[p.ID] = pAPI
// 	}
// 	return pAPIs, nil
// }
