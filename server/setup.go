package server

import (
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
)

func Setup() error {
	cs := cribbageServer{
		db: getDB(),
	}

	cs.Serve()

	return nil
}

func getDB() persistence.DB {
	return memory.New()
}
