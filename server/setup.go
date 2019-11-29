package server

import (
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
)

func Setup() error {
	cs := cribbageServer{
		db: getDB(),
	}
	// TODO is this where we should seed the DB with NPCs if they aren't
	// there already?

	cs.Serve()

	return nil
}

func getDB() persistence.DB {
	return memory.New()
}
