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
	err := seedNPCs()
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

func seedNPCs() error {
	ctx := context.Background()
	db, err := getDB(ctx)
	if err != nil {
		return err
	}

	npcIDs := []model.PlayerID{interaction.Dumb, interaction.Simple, interaction.Calc}
	for _, id := range npcIDs {
		p := model.Player{
			ID:   id,
			Name: string(id),
		}
		pm := interaction.New(id, interaction.Means{Mode: interaction.NPC})
		_, err := db.GetInteraction(id)
		if err != nil {
			if err != persistence.ErrInteractionNotFound {
				return err
			}
			err = db.SaveInteraction(pm)
			if err != nil {
				return err
			}
		}
		_, err = db.GetPlayer(p.ID)
		if err != nil {
			if err != persistence.ErrPlayerNotFound {
				return err
			}
			err = db.CreatePlayer(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
