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

	err := seedNPCs(cs)
	if err != nil {
		return err
	}

	cs.Serve()

	return nil
}

func seedNPCs(cs cribbageServer) error {
	npcIDs := []model.PlayerID{npc.Dumb, npc.Simple, npc.Calc}
	for _, id := range npcIDs {
		p, err := npc.NewNPCPlayer(id, cs.handleAction)
		if err != nil {
			return err
		}

		if _, err := cs.db.GetInteraction(p.ID()); err != nil {
			// TODO compare to err consts when those get pulled in
			if err.Error() == `does not have player` {
				return cs.db.SaveInteraction(p)
			}
			return err
		}
	}
	return nil
}

func getDB() persistence.DB {
	return memory.New()
}
