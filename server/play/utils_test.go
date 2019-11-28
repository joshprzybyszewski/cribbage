package play

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestPlayersToDealTo(t *testing.T) {
	alice := model.Player{
		ID:   model.PlayerID(1),
		Name: `alice`,
	}
	bob := model.Player{
		ID:   model.PlayerID(2),
		Name: `bob`,
	}
	charlie := model.Player{
		ID:   model.PlayerID(3),
		Name: `charlie`,
	}

	testCases := []struct{
		msg string
		g model.Game
		expPlayerIDs []model.PlayerID
	}{{
		msg: `two person game`,
		g: model.Game{
			Players: []model.Player{alice, bob},
			CurrentDealer: alice.ID,
		},
		expPlayerIDs: []model.PlayerID{bob.ID, alice.ID},
	}, {
		msg: `two person game, other dealer`,
		g: model.Game{
			Players: []model.Player{alice, bob},
			CurrentDealer: bob.ID,
		},
		expPlayerIDs: []model.PlayerID{alice.ID, bob.ID},
	}, {
		msg: `three person game`,
		g: model.Game{
			Players: []model.Player{alice, bob, charlie},
			CurrentDealer: alice.ID,
		},
		expPlayerIDs: []model.PlayerID{bob.ID, charlie.ID, alice.ID},
		}, {
			msg: `three person game, second dealer`,
			g: model.Game{
				Players: []model.Player{alice, bob, charlie},
				CurrentDealer: bob.ID,
			},
			expPlayerIDs: []model.PlayerID{charlie.ID, alice.ID, bob.ID},
			}, {
				msg: `three person game, third dealer`,
				g: model.Game{
					Players: []model.Player{alice, bob, charlie},
					CurrentDealer: charlie.ID,
				},
				expPlayerIDs: []model.PlayerID{ alice.ID, bob.ID, charlie.ID},
			
		}}

	for _, tc := range testCases {
		actPlayerIDs := playersToDealTo(&tc.g)
		assert.Equal(t, tc.expPlayerIDs, actPlayerIDs, tc.msg)
	}
}