package server

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction/npc"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
)

var (
	database = flag.String(`db`, `mongo`, `Set to the type of database to access`)
	dbURI    = flag.String(`dbURI`, ``, `The uri to the database. default empty string uses whatever localhost is`)
)

// Setup connects to a database and starts serving requests
func Setup() error {
	fmt.Printf("Using %s for persistence\n", *database)

	cs := cribbageServer{}
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

func getDB(ctx context.Context) (persistence.DB, error) {
	switch *database {
	case `mongo`:
		return mongodb.New(ctx, *dbURI)
	case `memory`:
		return memory.New(), nil
	}

	return nil, errors.New(`database type not supported`)
}
