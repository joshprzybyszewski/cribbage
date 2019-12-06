// +build !prod

package testutils

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func AliceAndBob() (a, b model.Player, am, bm *interaction.Mock, pAPIs map[model.PlayerID]interaction.Player) {
	alice := model.Player{
		ID:   model.PlayerID(`alice`),
		Name: `alice`,
	}
	bob := model.Player{
		ID:   model.PlayerID(`bob`),
		Name: `bob`,
	}
	aAPI := &interaction.Mock{}
	bAPI := &interaction.Mock{}
	abAPIs := map[model.PlayerID]interaction.Player{
		alice.ID: aAPI,
		bob.ID:   bAPI,
	}
	return alice, bob, aAPI, bAPI, abAPIs
}

func AliceBobCharlieDiane() (alice, bob, charlie, diane model.Player) {
	alice = model.Player{
		ID:   model.PlayerID(`a`),
		Name: `alice`,
	}
	bob = model.Player{
		ID:   model.PlayerID(`b`),
		Name: `bob`,
	}
	charlie = model.Player{
		ID:   model.PlayerID(`c`),
		Name: `charlie`,
	}
	diane = model.Player{
		ID:   model.PlayerID(`d`),
		Name: `diane`,
	}
	return alice, bob, charlie, diane
}
