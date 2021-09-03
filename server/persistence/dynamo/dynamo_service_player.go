//nolint:dupl
package dynamo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	playerServiceSortKeyPrefix string = `player`
	playerServiceGameSortKey   string = playerServiceSortKeyPrefix + `Game`
)

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	ctx context.Context

	svc *dynamodb.Client
}

func getPlayerService(
	ctx context.Context,
	svc *dynamodb.Client,
) (persistence.PlayerService, error) {

	return &playerService{
		ctx: ctx,
		svc: svc,
	}, nil
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	tableName := dbName
	pkName := `:pID`
	pk := string(id)
	skName := `:sk`
	sk := playerServiceSortKeyPrefix
	keyCondExpr := fmt.Sprintf("DDBid = %s and begins_with(spec, %s)", pkName, skName)
	qo, err := ps.svc.Query(ps.ctx, &dynamodb.QueryInput{
		TableName:              &tableName,
		KeyConditionExpression: &keyCondExpr,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			pkName: &types.AttributeValueMemberS{
				Value: pk,
			},
			skName: &types.AttributeValueMemberS{
				Value: sk,
			},
		},
	})
	if err != nil {
		return model.Player{}, err
	}
	// TODO check LastEvaluatedKey to know if we need to paginate the responses
	// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.Pagination
	ret := model.Player{
		ID:    id,
		Games: map[model.GameID]model.PlayerColor{},
	}
	foundName := false
	for _, item := range qo.Items {
		// TODO reduce the nesting...
		gID, color, ok := getGameIDAndColor(item[`spec`], item[`color`])
		if ok {
			ret.Games[gID] = color
			continue
		}
		name, ok := getPlayerName(item[`spec`], item[`name`])
		if !ok {
			return model.Player{}, errors.New(`got unexpected payload`)
		} else if foundName {
			return model.Player{}, errors.New(`found two names`)
		}
		foundName = true
		ret.Name = name
	}

	if !foundName {
		return model.Player{}, persistence.ErrPlayerNotFound
	}

	return ret, nil
}

func getGameIDAndColor(
	specAV, colorAV types.AttributeValue,
) (model.GameID, model.PlayerColor, bool) {
	specAVS, ok := specAV.(*types.AttributeValueMemberS)
	if !ok {
		return 0, 0, false
	}

	spec := specAVS.Value
	if len(spec) <= len(playerServiceGameSortKey) {
		return 0, 0, false
	}

	gID, err := strconv.Atoi(spec[len(playerServiceGameSortKey):])
	if err != nil {
		return 0, 0, false
	}

	colorAVN, ok := colorAV.(*types.AttributeValueMemberN)
	if !ok {
		return 0, 0, false
	}

	color, err := strconv.Atoi(colorAVN.Value)
	if err != nil {
		return 0, 0, false
	}

	return model.GameID(gID), model.PlayerColor(color), true
}

func getPlayerName(
	specAV, nameAV types.AttributeValue,
) (string, bool) {
	specAVS, ok := specAV.(*types.AttributeValueMemberS)
	if !ok {
		return ``, false
	}

	spec := specAVS.Value
	if spec != playerServiceSortKeyPrefix {
		return ``, false
	}

	nameAVS, ok := nameAV.(*types.AttributeValueMemberS)
	if !ok {
		return ``, false
	}
	return nameAVS.Value, true
}

func (ps *playerService) Create(p model.Player) error {
	if len(p.Games) > 0 {
		return errors.New(`you cannot create a player that is _already_ in games!`)
	}

	data := map[string]types.AttributeValue{
		`DDBid`: &types.AttributeValueMemberS{
			Value: string(p.ID),
		},
		`spec`: &types.AttributeValueMemberS{
			Value: playerServiceSortKeyPrefix,
		},
		`name`: &types.AttributeValueMemberS{
			Value: p.Name,
		},
	}

	// Use a conditional expression to only write items if this
	// <HASH:RANGE> tuple doesn't already exist.
	// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
	// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
	condExpr := `attribute_not_exists(DDBid) AND attribute_not_exists(spec)`

	fmt.Printf("playerService.Create data = %+v\n", data)

	_, err := ps.svc.PutItem(ps.ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(dbName),
		Item:                data,
		ConditionExpression: &condExpr,
	})
	if err != nil {
		switch err.(type) {
		case *types.ConditionalCheckFailedException:
			return persistence.ErrPlayerAlreadyExists
		}
		return err
	}

	return nil
}

type dynamoPlayerInGame struct {
	ID   string `json:"DDBid"`
	Spec string `json:"spec"`

	Color model.PlayerColor `json:"color"`
}

func (ps *playerService) BeginGame(gID model.GameID, players []model.Player) error {
	for _, p := range players {
		// TODO do these in separate goroutines (aka parallelize)
		err := ps.setPlayerGameColor(p.ID, gID, model.UnsetColor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *playerService) UpdateGameColor(
	pID model.PlayerID,
	gID model.GameID,
	color model.PlayerColor,
) error {

	p, err := ps.Get(pID)
	if err != nil {
		return err
	}

	if c, ok := p.Games[gID]; ok {
		if c != color {
			return errors.New(`mismatched player-games color`)
		}

		// Nothing to do; the player already knows its color
		return nil
	}

	return ps.setPlayerGameColor(pID, gID, color)
}

func (ps *playerService) setPlayerGameColor(
	pID model.PlayerID,
	gID model.GameID,
	color model.PlayerColor,
) error {
	data, err := attributevalue.MarshalMap(dynamoPlayerInGame{
		ID:    string(pID),
		Spec:  fmt.Sprintf("%s%d", playerServiceGameSortKey, gID),
		Color: color,
	})
	if err != nil {
		return err
	}

	_, err = ps.svc.PutItem(ps.ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item:      data,
	})
	return err
}
