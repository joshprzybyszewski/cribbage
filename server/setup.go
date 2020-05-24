package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mysql"
)

var (
	database = flag.String(`db`, `mysql`, `Set to the type of database to access. Options: "mysql", "mongo", "memory"`)
	dbURI    = flag.String(`dbURI`, ``, `The uri to the database. default empty string uses whatever localhost is`)

	dsnUser     = flag.String(`dsn_user`, `root`, `The DSN user for the MySQL DB`)
	dsnPassword = flag.String(`dsn_password`, ``, `The password for the user for the MySQL DB`)
	dsnHost     = flag.String(`dsn_host`, `127.0.0.1`, `The host for the MySQL DB`)
	dsnPort     = flag.Int(`dsn_port`, 3306, `The port for the MySQL DB`)
	dsnParams   = flag.String(`dsn_params`, ``, `The params for the MySQL DB`)
	mysqlDBName = flag.String(`mysql_db`, `cribbage`, `The name of the Database to connect to in mysql`)
)

// Setup connects to a database and starts serving requests
func Setup() error {
	fmt.Printf("Using %s for persistence\n", *database)

	db, err := getDB(context.Background())
	if err != nil {
		return err
	}
	cs := newCribbageServer(db)
	err = seedNPCs()
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
	case `mysql`:
		cfg := mysql.Config{
			DSNUser:      *dsnUser,
			DSNPassword:  *dsnPassword,
			DSNHost:      *dsnHost,
			DSNPort:      *dsnPort,
			DatabaseName: *mysqlDBName,
			DSNParams:    *dsnParams,
		}
		return mysql.New(ctx, cfg)
	case `memory`:
		return memory.New(), nil
	}

	return nil, fmt.Errorf(`db "%s" not supported. Currently supported: "mongo", "mysql", and "memory"`, *database)
}

func seedNPCs() error {
	ctx := context.Background()
	db, err := getDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Start()
	if err != nil {
		return err
	}
	defer commitOrRollback(db, &err)

	npcIDs := []model.PlayerID{interaction.Dumb, interaction.Simple, interaction.Calc}
	for _, id := range npcIDs {
		p := model.Player{
			ID:   id,
			Name: string(id),
		}
		_, err = db.GetPlayer(p.ID)
		if err != nil {
			if err == persistence.ErrPlayerNotFound {
				err = db.CreatePlayer(p)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		pm := interaction.New(id, interaction.Means{Mode: interaction.NPC})
		_, err = db.GetInteraction(id)
		if err != nil {
			if err == persistence.ErrInteractionNotFound {
				err = db.SaveInteraction(pm)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}
