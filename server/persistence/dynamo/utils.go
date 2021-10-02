package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func fullQuery(
	ctx context.Context,
	svc *dynamodb.Client,
	createQuery func() *dynamodb.QueryInput,
) ([]map[string]types.AttributeValue, error) {
	qi := createQuery()

	var items []map[string]types.AttributeValue
	for {
		qo, err := svc.Query(ctx, qi)
		if err != nil {
			return nil, err
		}
		items = append(items, qo.Items...)

		if len(qo.LastEvaluatedKey) == 0 {
			break
		}

		qi = createQuery()
		qi.ExclusiveStartKey = qo.LastEvaluatedKey
	}

	return items, nil
}
