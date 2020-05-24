package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestGetCribCards(t *testing.T) {
	// 0x35 = decimal 53 which is more than the num cards we have
	cards := getCribCards(int32(0x35353535))
	assert.Empty(t, cards)

	cards = getCribCards(int32(0x01020304))
	require.Len(t, cards, 4)
	assert.Equal(t, model.NewCardFromNumber(int(0x04)), cards[0])
	assert.Equal(t, model.NewCardFromNumber(int(0x03)), cards[1])
	assert.Equal(t, model.NewCardFromNumber(int(0x02)), cards[2])
	assert.Equal(t, model.NewCardFromNumber(int(0x01)), cards[3])

	cards = getCribCards(int32(0x01353504))
	require.Len(t, cards, 2)
	assert.Equal(t, model.NewCardFromNumber(int(0x04)), cards[0])
	assert.Equal(t, model.NewCardFromNumber(int(0x01)), cards[1])
}

func TestSerializeCribCards(t *testing.T) {
	serVal := serializeCribCards(nil)
	assert.Equal(t, int32(0x35353535), serVal)

	crib := []model.Card{
		model.NewCardFromNumber(int(0x05)),
		model.NewCardFromNumber(int(0x15)),
		model.NewCardFromNumber(int(0x08)),
		model.NewCardFromNumber(int(0x12)),
	}

	serVal = serializeCribCards(crib)
	assert.Equal(t, int32(0x12081505), serVal)

	crib = []model.Card{
		model.NewCardFromNumber(int(0x08)),
	}

	serVal = serializeCribCards(crib)
	assert.Equal(t, int32(0x35353508), serVal)
}
