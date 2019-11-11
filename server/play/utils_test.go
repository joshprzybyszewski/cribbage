package play

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandContains(t *testing.T) {
	assert.False(t, handContains(nil, model.NewCardFromString(`ah`)))
	assert.False(t, handContains([]model.Card{}, model.NewCardFromString(`ah`)))

	hand :=[]model.Card{
		model.NewCardFromString(`ah`),
		model.NewCardFromString(`jc`),
		model.NewCardFromString(`6c`),
		model.NewCardFromString(`8h`),
	}

	assert.True(t, handContains(hand, model.NewCardFromString(`ah`)))
	assert.True(t, handContains(hand, model.NewCardFromString(`jc`)))
	assert.False(t, handContains(hand, model.NewCardFromString(`ad`)))
	assert.False(t, handContains(hand, model.NewCardFromString(`7s`)))
	assert.False(t, handContains(hand, model.NewCardFromString(`jd`)))
}

func TestHasBeenPegged(t *testing.T) {
	assert.False(t, hasBeenPegged(nil, model.NewCardFromString(`ah`)))
	assert.False(t, hasBeenPegged([]model.PeggedCard{}, model.NewCardFromString(`ah`)))

	pegged :=[]model.PeggedCard{{
		Card: model.NewCardFromString(`ah`),
	}, {
		Card: model.NewCardFromString(`jc`),
	}}

	assert.True(t, hasBeenPegged(pegged, model.NewCardFromString(`ah`)))
	assert.True(t, hasBeenPegged(pegged, model.NewCardFromString(`jc`)))
	assert.False(t, hasBeenPegged(pegged, model.NewCardFromString(`ad`)))
	assert.False(t, hasBeenPegged(pegged, model.NewCardFromString(`7s`)))
	assert.False(t, hasBeenPegged(pegged, model.NewCardFromString(`jd`)))
}