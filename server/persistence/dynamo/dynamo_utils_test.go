package dynamo

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsConditionalError(t *testing.T) {
	assert.False(t, isConditionalError(nil))
	assert.False(t, isConditionalError(errors.New(`josh is cool`)))
	assert.True(t, isConditionalError(&types.ConditionalCheckFailedException{}))
	assert.True(t, isConditionalError(errors.New(`operation error DynamoDB: PutItem, https response error StatusCode: 400, RequestID: 68dd7a61-641c-4541-ac28-2de7556b8528, ConditionalCheckFailedException: `)))
}

func TestFullQuery(t *testing.T) {
	if testing.Short() {
		return
	}
	var err error
	ctx := context.Background()
	svc := getDynamoService(ctx, `http://localhost:18079`)

	hugeItemID := `hugeID` + rand.String(50)
	sk := `prefix`
	hugePayload := rand.String(390 * 1024)

	createData := func(i int) map[string]types.AttributeValue {
		payload := strconv.Itoa(i) + `suffix` + hugePayload
		return map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{
				Value: hugeItemID,
			},
			sortKey: &types.AttributeValueMemberS{
				Value: sk + strconv.Itoa(i),
			},
			`dummyData`: &types.AttributeValueMemberS{
				Value: payload,
			},
		}
	}

	numGen := 5

	for i := 0; i < numGen; i++ {
		data := createData(i)

		_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(dbName),
			Item:      data,
		})
		require.NoError(t, err)
	}

	pkName := `:pID`
	skName := `:sk`
	keyCondExpr := getConditionExpression(equalsID, pkName, hasPrefix, skName)

	createQuery := func() *dynamodb.QueryInput {
		return &dynamodb.QueryInput{
			TableName:              aws.String(dbName),
			KeyConditionExpression: &keyCondExpr,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				pkName: &types.AttributeValueMemberS{
					Value: hugeItemID,
				},
				skName: &types.AttributeValueMemberS{
					Value: sk,
				},
			},
		}
	}
	items, err := fullQuery(ctx, svc, createQuery)
	require.NoError(t, err)

	if len(items) != numGen {
		t.Errorf("did not have %d items, but had %d", numGen, len(items))
	}
}
