package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

func getPlayerAPIs(db persistence.DB, players []model.Player) (map[model.PlayerID]interaction.Player, error) {
	pAPIs := make(map[model.PlayerID]interaction.Player, len(players))
	for _, p := range players {
		pm, err := db.GetInteraction(p.ID)
		if err != nil {
			return nil, err
		}
		pAPI, err := interaction.FromPlayerMeans(pm)
		if err != nil {
			return nil, err
		}
		pAPIs[p.ID] = pAPI
	}
	return pAPIs, nil
}
