package server

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction/npc"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
)

// Setup connects to a database and starts serving requests
func Setup() error {
	cs := cribbageServer{
		db: getDB(),
	}

	for _, id := range getNPCIDs() {
		npc, err := npc.NewNPCPlayer(id, cs.handleAction)
		if err != nil {
			return err
		}

		if _, err := cs.db.GetInteraction(npc.ID()); err != nil {
			// TODO we should probably not use the error this way...
			if err.Error() == `does not have player` {
				return cs.db.SaveInteraction(npc)
			}
			return err
		}
	}

	cs.Serve()

	return nil
}

func getNPCIDs() []model.PlayerID {
	return []model.PlayerID{`dumbNPC`, `simpleNPC`, `calculatedNPC`}
}

func getDB() persistence.DB {
	return memory.New()
}
