package dynamo

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	dbName       = `cribbage`
	partitionKey = `DDBid`
	sortKey      = `spec`
)

func getSortKeyPrefix(service interface{}) string {
	// all of these need to be different, because they are the
	// start of the sort key. we are partitioning our dynamo table usage
	// such that each service has the same prefix:#
	switch service.(type) {
	case *gameService:
		return `game`
	case *interactionService:
		return `interaction`
	case *playerService:
		return `player`
	}

	return `garbage`
}

type hasPrefix struct {
	pkName string
	skName string
}

func (hp hasPrefix) conditionExpression() *string {
	var sb strings.Builder

	sb.WriteString(partitionKey)
	sb.WriteString(`=`)
	sb.WriteString(hp.pkName)

	sb.WriteString(` and `)

	sb.WriteString(`begins_with(`)
	sb.WriteString(sortKey)
	sb.WriteString(`,`)
	sb.WriteString(hp.skName)
	sb.WriteString(`)`)

	return aws.String(sb.String())
}

type notExists struct{}

func (notExists) conditionExpression() *string {
	var sb strings.Builder

	sb.WriteString(`attribute_not_exists(`)
	sb.WriteString(partitionKey)
	sb.WriteString(`)`)

	sb.WriteString(` and `)

	sb.WriteString(`attribute_not_exists(`)
	sb.WriteString(sortKey)
	sb.WriteString(`)`)

	return aws.String(sb.String())
}

type queryInputParams struct {
	partitionKey     string
	partitionKeyName string
	sortKey          string
	sortKeyName      string

	keyConditionExpression *string
}

func getQueryInputParams(
	pk, pkName,
	sk, skName string,
	cond *string,
) queryInputParams {
	return queryInputParams{
		partitionKey:           pk,
		partitionKeyName:       pkName,
		sortKey:                sk,
		sortKeyName:            skName,
		keyConditionExpression: cond,
	}
}

func newQueryInputFactory(
	params queryInputParams,
) func() *dynamodb.QueryInput {
	return func() *dynamodb.QueryInput {
		return &dynamodb.QueryInput{
			TableName:              aws.String(dbName),
			KeyConditionExpression: params.keyConditionExpression,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				params.partitionKeyName: &types.AttributeValueMemberS{
					Value: params.partitionKey,
				},
				params.sortKeyName: &types.AttributeValueMemberS{
					Value: params.sortKey,
				},
			},
		}
	}
}
