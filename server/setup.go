package server

import (
	"errors"
	"flag"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
)

var (
	database = flag.String(`db`, `memory`, `Set to the type of database to access`)
)

func Setup() error {
	db, err := getDB()
	if err != nil {
		return err
	}

	cs := cribbageServer{
		db: db,
	}

	cs.Serve()

	return nil
}

func getDB() (persistence.DB, error) {
	switch *database {
	case `mongo`:
		return mongodb.New(`mongodb://localhost:27017`)
	case `memory`:
		return memory.New(), nil
	}

	return nil, errors.New(`database type not supported`)
}
