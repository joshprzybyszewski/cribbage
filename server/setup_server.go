package server

import (
	"context"
	"log"
)

// Setup connects to a database and starts serving requests
func Setup() error {
	loadConfig()
	log.Printf("Using %s for persistence\n", *database)

	ctx := context.Background()
	dbFactory, err := getDBFactory(ctx, factoryConfig{
		canRunCreateStmts: true,
	})
	if err != nil {
		return err
	}
	cs := newCribbageServer(dbFactory)
	if err := seedNPCs(ctx, dbFactory); err != nil {
		return err
	}

	return cs.Serve()
}
