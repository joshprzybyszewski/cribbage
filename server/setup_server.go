//+build !lambda

package server

import (
	"context"
	"log"
)

// Setup connects to a database and starts serving requests
func Setup() error {
	loadVarsFromINI()
	log.Printf("Using %s for persistence\n", *database)

	ctx := context.Background()
	dbFactory, err := getDBFactory(ctx, factoryConfig{
		canRunCreateStmts: true,
	})
	if err != nil {
		return err
	}
	cs := newCribbageServer(dbFactory)
	err = seedNPCs(ctx, dbFactory)
	if err != nil {
		return err
	}
	cs.Serve()

	return nil
}
