package dynamo

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func TestGetSortKeyPrefix(t *testing.T) {
	testCases := []struct {
		service   interface{}
		expPrefix string
	}{{
		service:   (*gameService)(nil),
		expPrefix: `game`,
	}, {
		service:   (*interactionService)(nil),
		expPrefix: `interaction`,
	}, {
		service:   (*playerService)(nil),
		expPrefix: `player`,
	}, {
		service:   (*model.Game)(nil),
		expPrefix: `garbage`,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expPrefix, getSortKeyPrefix(tc.service))
	}
}

func TestGetConditionExpression(t *testing.T) {
	testCases := []struct {
		pkt, skt condExprType
		pk, sk   string
		exp      string
	}{{
		pkt: equalsID,
		pk:  `:pkAttrName`,
		skt: hasPrefix,
		sk:  `:skAttrName`,
		exp: `DDBid=:pkAttrName and begins_with(spec,:skAttrName)`,
	}, {
		pkt: notExists,
		pk:  `:pkAttrName`,
		skt: notExists,
		sk:  `:skAttrName`,
		exp: `attribute_not_exists(DDBid) and attribute_not_exists(spec)`,
	}, {
		pkt: equalsID,
		pk:  `:pkAttrName`,
		skt: none,
		sk:  `:skAttrName`,
		exp: `DDBid=:pkAttrName`,
	}, {
		pkt: hasPrefix,
		pk:  `:pkAttrName`,
		skt: equalsID,
		sk:  `:skAttrName`,
		exp: `unsupported pkType and unsupported skType`,
	}}

	for _, tc := range testCases {
		act := getConditionExpression(tc.pkt, tc.pk, tc.skt, tc.sk)
		assert.Equal(t, tc.exp, act)
	}
}
