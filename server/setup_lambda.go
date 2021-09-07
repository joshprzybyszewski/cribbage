//+build lambda

package server

import (
	"context"
	"fmt"
	"log"

	"github.com/apex/gateway"
	"github.com/rakyll/globalconf"
)

// Setup for lambda does not include static content, but serves all the API routes with persistence specified by
// the CRIBBAGE_DB environment variable
func Setup() error {
	cs, err := newServer()
	if err != nil {
		return err
	}
	return gateway.ListenAndServe(`:8080`, cs.NewRouter())
}

func newServer() (*cribbageServer, error) {
	loadFlagsFromEnvVars()
	log.Printf("Using %s for persistence\n", *database)

	dbFactory, err := getDBFactory(context.Background(), factoryConfig{
		canRunCreateStmts: true,
	})
	if err != nil {
		return nil, err
	}
	return newCribbageServer(dbFactory), nil
}

// for lambda, we want to only use env vars. Load globalconf without setting the file
func loadFlagsFromEnvVars() {
	conf, err := globalconf.NewWithOptions(&globalconf.Options{
		EnvPrefix: `CRIBBAGE_`,
	})
	if err != nil {
		panic(fmt.Sprintf("globalconf.NewWithOptions error: %+v", err))
	}

	conf.ParseAll()
}
