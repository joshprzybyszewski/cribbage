package interaction

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func Test_buildNPCAction(t *testing.T) {
	testCases := []struct {
		desc string
		npc  NPC
		b    model.Blocker
		g    model.Game
		exp  model.PlayerAction
	}{
		// TODO: Add test cases.
	}
	for _, tc := range testCases {
		a := buildNPCAction(tc.npc, tc.b, tc.g)
		assert.Equal(t, tc.exp, a)
	}
}
