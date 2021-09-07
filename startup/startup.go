//+build !lambda

package startup

import "github.com/joshprzybyszewski/cribbage/server"

func PlayServer() error {
	return server.Setup()
}
