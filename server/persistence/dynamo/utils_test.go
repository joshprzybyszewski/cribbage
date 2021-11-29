package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
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
		get interface{ conditionExpression() *string }
		exp *string
	}{{
		get: hasPrefix{
			pkName: `:pkAttrName`,
			skName: `:skAttrName`,
		},
		exp: aws.String(`cribbageID=:pkAttrName and begins_with(spec,:skAttrName)`),
	}, {
		get: notExists{},
		exp: aws.String(`attribute_not_exists(cribbageID) and attribute_not_exists(spec)`),
	}}

	for _, tc := range testCases {
		act := tc.get.conditionExpression()
		assert.Equal(t, tc.exp, act)
	}
}

func TestNewQueryInputFactory(t *testing.T) {
	testCases := []struct {
		pk, pkName string
		sk, skName string
		cond       *string
		exp        *dynamodb.QueryInput
	}{{
		pk:     `is_cool`,
		pkName: `:josh`,
		sk:     `so_cool`,
		skName: `:jp`,
		cond:   aws.String(`attribute_not_exists(:josh) and attribute_not_exists(:jp)`),
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
		createQuery := newQueryInputFactory(getQueryInputParams(
			tc.pk, tc.pkName,
			tc.sk, tc.skName,
			tc.cond,
		))
		require.NotNil(t, createQuery)
		actQuery := createQuery()
		assert.Equal(t, tc.exp, actQuery)
	}
}
