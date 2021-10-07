package dynamo

import "strings"

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

type condExprType uint8

const (
	none      condExprType = 0
	equalsID  condExprType = 1
	hasPrefix condExprType = 2
	notExists condExprType = 3
)

func getConditionExpression(
	pkType condExprType,
	pk string,
	skType condExprType,
	sk string,
) string {
	var sb strings.Builder

	switch pkType {
	case equalsID:
		sb.WriteString(partitionKey)
		sb.WriteString(`=`)
		sb.WriteString(pk)
	case notExists:
		sb.WriteString(`attribute_not_exists(`)
		sb.WriteString(partitionKey)
		sb.WriteString(`)`)
	default:
		sb.WriteString(`unsupported pkType`)
	}

	if skType == none {
		return sb.String()
	}

	sb.WriteString(` and `)

	switch skType {
	case hasPrefix:
		sb.WriteString(`begins_with(`)
		sb.WriteString(sortKey)
		sb.WriteString(`,`)
		sb.WriteString(sk)
		sb.WriteString(`)`)
	case notExists:
		sb.WriteString(`attribute_not_exists(`)
		sb.WriteString(sortKey)
		sb.WriteString(`)`)
	default:
		sb.WriteString(`unsupported skType`)
	}

	return sb.String()
}
