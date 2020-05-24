package persistence

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateLatestActionBelongs(t *testing.T) {
	mg := model.Game{}

	assert.NoError(t, ValidateLatestActionBelongs(mg))

	gID := model.GameID(1212)
	notGID := model.GameID(420)
	mg.ID = gID
	p1ID := model.PlayerID(`123`)
	p2ID := model.PlayerID(`456`)
	noPID := model.PlayerID(`789`)
	mg.Players = []model.Player{{
		ID: p1ID,
	}, {
		ID: p2ID,
	}}
	assert.NoError(t, ValidateLatestActionBelongs(mg))

	mg.Actions = append(mg.Actions, model.PlayerAction{
		ID:     p1ID,
		GameID: gID,
	})
	assert.NoError(t, ValidateLatestActionBelongs(mg))

	mg.Actions = append(mg.Actions, model.PlayerAction{
		ID:     p1ID,
		GameID: notGID,
	})
	assert.Error(t, ValidateLatestActionBelongs(mg))

	mg.Actions = append(mg.Actions, model.PlayerAction{
		ID:     noPID,
		GameID: gID,
	})
	assert.Error(t, ValidateLatestActionBelongs(mg))

	mg.Actions = append(mg.Actions, model.PlayerAction{
		ID:     p2ID,
		GameID: gID,
	})
	assert.NoError(t, ValidateLatestActionBelongs(mg))
}
