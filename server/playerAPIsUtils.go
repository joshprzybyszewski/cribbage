package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var actionHandler = &npcActionHandler{}

func getPlayerAPIs(db persistence.DB, players []model.Player) (map[model.PlayerID]interaction.Player, error) {
	pAPIs := make(map[model.PlayerID]interaction.Player, len(players))
	for _, p := range players {
		var pAPI interaction.Player
		pm, err := db.GetInteraction(p.ID)

		for i, m := range pm.Interactions {
			if m.Mode == interaction.NPC {
				m.Info = actionHandler
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
