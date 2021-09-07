//+build lambda

package startup

import (
	"github.com/apex/gateway"
	"github.com/joshprzybyszewski/cribbage/server"
)

func PlayServer() error {
	cs, err := server.NewServer()
	if err != nil {
		return err
	}
	return gateway.ListenAndServe(`:8080`, cs.NewRouter())
}
