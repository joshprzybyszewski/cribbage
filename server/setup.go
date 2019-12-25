package server

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
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
	ctx := context.Background()
	db, err := getDB(ctx)
	if err != nil {
		return err
	}
	err = seedNPCs(db)
	if err != nil {
		return err
	}
	cs.Serve()

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

func seedNPCs(db persistence.DB) error {
	npcIDs := []model.PlayerID{interaction.Dumb, interaction.Simple, interaction.Calc}
	for _, id := range npcIDs {
		p := model.Player{
			ID:   id,
			Name: string(id),
		}
		if _, err := db.GetPlayer(p.ID); err != nil {
			if err == persistence.ErrPlayerNotFound {
				return db.CreatePlayer(p)
			}
			return err
		}
	}
	return nil
}
