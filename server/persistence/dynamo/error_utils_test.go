package dynamo

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestIsConditionalError(t *testing.T) {
	assert.False(t, isConditionalError(nil))
	assert.False(t, isConditionalError(errors.New(`josh is cool`)))
	assert.True(t, isConditionalError(&types.ConditionalCheckFailedException{}))
	assert.True(t, isConditionalError(errors.New(`operation error DynamoDB: PutItem, https response error StatusCode: 400, RequestID: 68dd7a61-641c-4541-ac28-2de7556b8528, ConditionalCheckFailedException: `)))
}
