package jsonutils

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestUnmarshalPlayerAction(t *testing.T) {
	testCases := []struct {
		msg string
		pa  model.PlayerAction
	}{{
		msg: `deal`,
		pa: model.PlayerAction{
			GameID:    model.GameID(42),
			ID:        model.PlayerID(`alice`),
			Overcomes: model.DealCards,
			Action: model.DealAction{
				NumShuffles: 543,
			},
			TimeStamp: time.Now(),
		},
	}, {
		msg: `crib`,
		pa: model.PlayerAction{
			GameID:    model.GameID(45),
			ID:        model.PlayerID(`bob`),
			Overcomes: model.CribCard,
			Action: model.BuildCribAction{
				Cards: []model.Card{
					model.NewCardFromString(`jh`),
					model.NewCardFromString(`5d`),
				},
			},
			TimeStamp: time.Now(),
		},
	}, {
		msg: `cut`,
		pa: model.PlayerAction{
			GameID:    model.GameID(312),
			ID:        model.PlayerID(`charlie`),
			Overcomes: model.CutCard,
			Action: model.CutDeckAction{
				Percentage: 0.12345,
			},
		},
	}, {
		msg: `peg`,
		pa: model.PlayerAction{
			GameID:    model.GameID(999),
			ID:        model.PlayerID(`diane`),
			Overcomes: model.PegCard,
			Action: model.PegAction{
				Card: model.NewCardFromString(`jh`),
			},
			TimeStamp: time.Now(),
		},
	}, {
		msg: `peg saygo`,
		pa: model.PlayerAction{
			GameID:    model.GameID(1),
			ID:        model.PlayerID(`edward`),
			Overcomes: model.PegCard,
			Action: model.PegAction{
				SayGo: true,
			},
			TimeStamp: time.Now(),
		},
	}, {
		msg: `count hand`,
		pa: model.PlayerAction{
			GameID:    model.GameID(54),
			ID:        model.PlayerID(`frances`),
			Overcomes: model.CountHand,
			Action: model.CountHandAction{
				Pts: 29,
			},
			TimeStamp: time.Now(),
		},
	}, {
		msg: `count crib`,
		pa: model.PlayerAction{
			GameID:    model.GameID(3),
			ID:        model.PlayerID(`george`),
			Overcomes: model.CountCrib,
			Action: model.CountCribAction{
				Pts: 12,
			},
			TimeStamp: time.Now(),
		},
	}}

	for _, tc := range testCases {
		tc.pa.SetTimeStamp(time.Now())
		paCopy := tc.pa
		b, err := json.Marshal(tc.pa)
		require.NoError(t, err, tc.msg)
		actPA, err := UnmarshalPlayerAction(b)
		require.NoError(t, err, tc.msg)
		assert.Equal(t, paCopy, actPA, tc.msg)
	}
}
