// +build !prod

package testutils

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

func aliceAndBob() (a, b model.Player) {
	alice := model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `alice`,
	}
	bob := model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `bob`,
	}
	return alice, bob
}

func AliceAndBob() (a, b model.Player, am, bm *interaction.Mock, pAPIs map[model.PlayerID]interaction.Player) {
	alice, bob := aliceAndBob()
	aAPI := &interaction.Mock{}
	bAPI := &interaction.Mock{}
	abAPIs := map[model.PlayerID]interaction.Player{
		alice.ID: aAPI,
		bob.ID:   bAPI,
	}
	return alice, bob, aAPI, bAPI, abAPIs
}

func EmptyAliceAndBob() (a, b model.Player, pAPIs map[model.PlayerID]interaction.Player) {
	alice, bob := aliceAndBob()
	aAPI := interaction.Empty(alice.ID)
	bAPI := interaction.Empty(bob.ID)
	abAPIs := map[model.PlayerID]interaction.Player{
		alice.ID: aAPI,
		bob.ID:   bAPI,
	}
	return alice, bob, abAPIs
}

func AliceBobCharlieDiane() (alice, bob, charlie, diane model.Player) {
	alice = model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `alice`,
	}
	bob = model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `bob`,
	}
	charlie = model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `charlie`,
	}
	diane = model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `diane`,
	}
	return alice, bob, charlie, diane
}
