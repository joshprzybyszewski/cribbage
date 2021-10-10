package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetPkSkCondCreateQuery(t *testing.T) {
	testCases := []struct {
		pk, pkName string
		sk, skName string
		cond       string
		exp        *dynamodb.QueryInput
	}{{
		pk:     `is_cool`,
		pkName: `:josh`,
		sk:     `so_cool`,
		skName: `:jp`,
		cond:   `attribute_not_exists(:josh) and attribute_not_exists(:jp)`,
		exp: &dynamodb.QueryInput{
			TableName:              aws.String(dbName),
			KeyConditionExpression: aws.String(`attribute_not_exists(:josh) and attribute_not_exists(:jp)`),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				`:josh`: &types.AttributeValueMemberS{
					Value: `is_cool`,
				},
				`:jp`: &types.AttributeValueMemberS{
					Value: `so_cool`,
				},
			},
		},
	}}

	for _, tc := range testCases {
		createQuery := getPkSkCondCreateQuery(
			tc.pk, tc.pkName,
			tc.sk, tc.skName,
			tc.cond,
		)
		require.NotNil(t, createQuery)
		actQuery := createQuery()
		assert.Equal(t, tc.exp, actQuery)
	}
}
