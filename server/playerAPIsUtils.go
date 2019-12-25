package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

func getPlayerAPIs(db persistence.DB, players []model.Player) (map[model.PlayerID]interaction.Player, error) {
	pAPIs := make(map[model.PlayerID]interaction.Player, len(players))
	for _, p := range players {
		var pAPI interaction.Player
		pm, err := db.GetInteraction(p.ID)

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
			// TODO assign HandleAction callback here if the character is an NPC
			if pAPI, ok := pAPI.(*interaction.NPCPlayer); ok {
				pAPI.HandleActionCallback = HandleAction
			}
			pAPI.ID()
		}

		pAPIs[p.ID] = pAPI
	}
	return pAPIs, nil
}
