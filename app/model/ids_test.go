package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGameID(t *testing.T) {
	for i := 0; i < 100; i++ {
		gID := NewGameID()
		require.NotEqual(t, InvalidGameID, gID)
	}
}

func TestIsValidPlayerID(t *testing.T) {
	testCases := []struct {
		msg     string
		input   string
		isValid bool
	}{{
		msg:     `normal stuff`,
		input:   `normalStuff`,
		isValid: true,
	}, {
		msg:     `has numbers`,
		input:   `w1thnumb3r5`,
		isValid: true,
	}, {
		msg:     `has underscores`,
		input:   `has_under_scores`,
		isValid: true,
	}, {
		msg:     `caps`,
		input:   `hAsCaPiTaLlEtTeRs`,
		isValid: true,
	}, {
		msg:     `has dashes`,
		input:   `has-dashes-yanno`,
		isValid: false,
	}, {
		msg:     `spaces`,
		input:   `has spaces dude`,
		isValid: false,
	}, {
		msg:     `special chars`,
		input:   `what!`,
		isValid: false,
	}, {
		msg:     `empty string`,
		input:   ``,
		isValid: false,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.isValid, IsValidPlayerID(PlayerID(tc.input)))
	}
}
