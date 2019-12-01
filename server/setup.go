package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
)

// Setup connects to a database and starts serving requests
func Setup() error {
	cs := cribbageServer{
		db: getDB(),
	}

	npcTypes := []interaction.NPC{
		interaction.Dumb, interaction.Simple, interaction.Calculated,
	}
	for _, npcType := range npcTypes {
		npc := interaction.NewNPCPlayer(npcType)
		if _, err := cs.db.GetPlayer(npc.ID()); err != nil {
			// TODO we should probably not use the error this way...
			if err.Error() == `does not have player` {
				// TODO how should we store NPCs in the database? Or should
				// we?
				p := model.Player{
					ID:   npc.ID(),
					Name: ``,
				}
				return cs.db.CreatePlayer(p)
			}
			return err
		}
	}

	cs.Serve()

	return nil
}

func getDB() persistence.DB {
	return memory.New()
}
