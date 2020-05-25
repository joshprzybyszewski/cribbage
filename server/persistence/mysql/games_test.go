package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
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

func TestSerializeAndDeserializeAllTheThings(t *testing.T) {
	alice, bob, _, _, _ := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.CountCrib},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.CribCounting,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: {
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`9s`),
				model.NewCardFromString(`10s`),
			},
			bob.ID: {
				model.NewCardFromString(`7c`),
				model.NewCardFromString(`8c`),
				model.NewCardFromString(`9c`),
				model.NewCardFromString(`10c`),
			},
		},
		CutCard: model.NewCardFromString(`7h`),
		Crib: []model.Card{
			model.NewCardFromString(`7d`),
			model.NewCardFromString(`8d`),
			model.NewCardFromString(`9d`),
			model.NewCardFromString(`10d`),
		},
		PeggedCards: []model.PeggedCard{
			model.NewPeggedCard(bob.ID, model.NewCardFromString(`7c`), 0),
			model.NewPeggedCard(alice.ID, model.NewCardFromString(`7s`), 0),
			model.NewPeggedCard(bob.ID, model.NewCardFromString(`8c`), 0),
			model.NewPeggedCard(alice.ID, model.NewCardFromString(`8s`), 0),
			model.NewPeggedCard(bob.ID, model.NewCardFromString(`9c`), 0),
			model.NewPeggedCard(alice.ID, model.NewCardFromString(`9s`), 0),
			model.NewPeggedCard(bob.ID, model.NewCardFromString(`10c`), 0),
			model.NewPeggedCard(alice.ID, model.NewCardFromString(`10s`), 0),
		},
	}

	cribCopy := make([]model.Card, len(g.Crib))
	for i := range cribCopy {
		cribCopy[i] = g.Crib[i]
	}
	assert.Equal(t, cribCopy, getCribCards(serializeCribCards(g.Crib)))

	bpCpy := make(map[model.PlayerID]model.Blocker, len(g.BlockingPlayers))
	for k, v := range g.BlockingPlayers {
		bpCpy[k] = v
	}
	bpSer, err := serializeBlockingPlayers(g.BlockingPlayers)
	require.NoError(t, err)
	actBlocking, err := getBlockingPlayers(bpSer)
	require.NoError(t, err)
	assert.Equal(t, bpCpy, actBlocking)

	handCpy := make(map[model.PlayerID][]model.Card, len(g.Hands))
	for k, pHand := range g.Hands {
		pHandCpy := make([]model.Card, len(pHand))
		for i := range pHandCpy {
			pHandCpy[i] = pHand[i]
		}
		handCpy[k] = pHandCpy
	}
	hSer, err := serializeHands(g.Hands)
	require.NoError(t, err)
	actHands, err := getHands(hSer)
	require.NoError(t, err)
	assert.Equal(t, handCpy, actHands)

	peggedCpy := make([]model.PeggedCard, len(g.PeggedCards))
	for i := range peggedCpy {
		peggedCpy[i] = g.PeggedCards[i]
	}
	pcSer, err := serializePeggedCards(g.PeggedCards)
	require.NoError(t, err)
	actPC, err := getPeggedCards(pcSer)
	require.NoError(t, err)
	assert.Equal(t, peggedCpy, actPC)
}
