//nolint:dupl
package dynamo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
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

func (ps *playerService) getGameSortKey() string {
	return getSortKeyPrefix(ps) + `Game`
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	ret := model.Player{
		ID:    id,
		Games: map[model.GameID]model.PlayerColor{},
	}

	pkName := `:pID`
	pk := string(id)
	skName := `:sk`
	sk := getSortKeyPrefix(ps)
	keyCondExpr := getConditionExpression(equalsID, pkName, hasPrefix, skName)

	createQuery := func() *dynamodb.QueryInput {
		return &dynamodb.QueryInput{
			TableName:              aws.String(dbName),
			KeyConditionExpression: &keyCondExpr,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				pkName: &types.AttributeValueMemberS{
					Value: pk,
				},
				skName: &types.AttributeValueMemberS{
					Value: sk,
				},
			},
		}
	}
	qi := createQuery()

	for {
		qo, err := ps.svc.Query(ps.ctx, qi)
		if err != nil {
			return model.Player{}, err
		}

		err = ps.populatePlayerFromItems(&ret, qo.Items)
		// check LastEvaluatedKey to know if we need to paginate the responses
		// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.Pagination
		if len(qo.LastEvaluatedKey) == 0 {
			break
		}

		qi = createQuery()
		qi.ExclusiveStartKey = qo.LastEvaluatedKey
	}

	if ret.Name == `` {
		// This player _must_ have had at least a name stored, otherwise
		// we done messed up
		return model.Player{}, persistence.ErrPlayerNotFound
	}

	return ret, nil
}

func (ps *playerService) populatePlayerFromItems(
	p *model.Player,
	items []map[string]types.AttributeValue,
) error {
	for _, item := range items {
		if colorAV, ok := item[ps.getColorKey()]; ok {
			gID, color, err := ps.getGameIDAndColor(item[sortKey], colorAV)
			if err != nil {
				return err
			}

			p.Games[gID] = color
			continue
		}

		name, ok := ps.getPlayerName(item[sortKey], item[`name`])
		if !ok {
			return errors.New(`got unexpected payload`)
		} else if p.Name != `` {
			return errors.New(`found two names`)
		}

		p.Name = name
	}

	return nil
}

func (ps *playerService) getSpecForPlayerGameColor(
	gID model.GameID,
) string {
	return fmt.Sprintf(`%s%d`, ps.getGameSortKey(), gID)
}

func (ps *playerService) getPlayerGameColorFromSpec(
	spec string,
) (model.GameID, error) {
	if len(spec) <= len(ps.getGameSortKey()) {
		return 0, errors.New(`too short`)
	}

	gID, err := strconv.Atoi(spec[len(ps.getGameSortKey()):])
	if err != nil {
		return 0, err
	}

	return model.GameID(gID), nil
}

func (ps *playerService) getGameIDAndColor(
	specAV, colorAV types.AttributeValue,
) (model.GameID, model.PlayerColor, error) {
	specAVS, ok := specAV.(*types.AttributeValueMemberS)
	if !ok {
		return 0, 0, errors.New(`spec wrong type`)
	}

	gID, err := ps.getPlayerGameColorFromSpec(specAVS.Value)
	if err != nil {
		return 0, 0, err
	}

	colorAVS, ok := colorAV.(*types.AttributeValueMemberS)
	if !ok {
		return 0, 0, errors.New(`color wrong type`)
	}
	pc := model.NewPlayerColorFromString(colorAVS.Value)

	return gID, pc, nil
}

func (ps *playerService) getPlayerName(
	specAV, nameAV types.AttributeValue,
) (string, bool) {
	specAVS, ok := specAV.(*types.AttributeValueMemberS)
	if !ok {
		return ``, false
	}

	spec := specAVS.Value
	if spec != getSortKeyPrefix(ps) {
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
		partitionKey: &types.AttributeValueMemberS{
			Value: string(p.ID),
		},
		sortKey: &types.AttributeValueMemberS{
			Value: getSortKeyPrefix(ps),
		},
		`name`: &types.AttributeValueMemberS{
			Value: p.Name,
		},
	}

	// Use a conditional expression to only write items if this
	// <HASH:RANGE> tuple doesn't already exist.
	// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
	// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
	condExpr := getConditionExpression(notExists, ``, notExists, ``)

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

	if c, ok := p.Games[gID]; ok && c != model.UnsetColor {
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
	_, err := ps.svc.PutItem(ps.ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{
				Value: string(pID),
			},
			sortKey: &types.AttributeValueMemberS{
				Value: ps.getSpecForPlayerGameColor(gID),
			},
			ps.getColorKey(): &types.AttributeValueMemberS{
				Value: color.String(),
			},
		},
	})
	return err
}

func (ps *playerService) getColorKey() string {
	return `color`
}
