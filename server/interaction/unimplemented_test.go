package interaction

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func TestUnimplemented(t *testing.T) {
	myID := model.PlayerID(`voltron`)

	p := newUnimplemented(myID)

	assert.Equal(t, myID, p.ID())
	assert.Nil(t, p.NotifyBlocking(model.DealCards, model.Game{}, ``))
	assert.Nil(t, p.NotifyMessage(model.Game{}, ``))
	assert.Nil(t, p.NotifyScoreUpdate(model.Game{}, ``, ``))
}
