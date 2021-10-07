package dynamo

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func isConditionalError(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *types.ConditionalCheckFailedException:
		return true
	}

	return strings.Contains(
		err.Error(),
		(*types.ConditionalCheckFailedException)(nil).ErrorCode(),
	)
}
