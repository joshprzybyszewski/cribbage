package strategy

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func TestPegHighestCardNow(t *testing.T) {
	tests := []struct {
		desc     string
		hand     []model.Card
		prevPegs []model.PeggedCard
		curPeg   int
		expGo    bool
		expCard  model.Card
	}{{
		desc: `test no cards pegged yet`,
		hand: []model.Card{
			model.NewCardFromString(`ah`),
			model.NewCardFromString(`2h`),
			model.NewCardFromString(`3h`),
			model.NewCardFromString(`4h`),
		},
		prevPegs: make([]model.PeggedCard, 0),
		curPeg:   0,
		expGo:    false,
		expCard:  model.NewCardFromString(`ah`),
	}}
	for _, tc := range tests {
		c, sayGo := PegHighestCardNow(tc.hand, tc.prevPegs, tc.curPeg)
		assert.Equal(t, tc.expCard, c)
		assert.Equal(t, tc.expGo, sayGo)
	}
}
