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

const (
	playerColorAttributeName = `color`
	playerNameAttributeName  = `name`
)

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	ctx context.Context

	svc *dynamodb.Client
}

func newPlayerService(
	ctx context.Context,
	svc *dynamodb.Client,
) persistence.PlayerService {
	return &playerService{
		ctx: ctx,
		svc: svc,
	}
}

func (ps *playerService) getGameSortKey() string {
	return getSortKeyPrefix(ps) + `Game`
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	pkName := `:pID`
	pk := string(id)
	skName := `:sk`
	sk := getSortKeyPrefix(ps)
	keyCondExpr := getConditionExpression(equalsID, pkName, hasPrefix, skName)

	createQuery := newQueryInputFactory(getQueryInputParams(
		pk, pkName, sk, skName, keyCondExpr,
	))
	items, err := fullQuery(ps.ctx, ps.svc, createQuery)
	if err != nil {
		return model.Player{}, err
	}

	return ps.buildPlayerFromItems(id, items)
}

func (ps *playerService) buildPlayerFromItems(
	id model.PlayerID,
	items []map[string]types.AttributeValue,
) (model.Player, error) {

	p := model.Player{
		ID:    id,
		Games: map[model.GameID]model.PlayerColor{},
	}

	for _, item := range items {
		if colorAV, ok := item[ps.getColorKey()]; ok {
			gID, color, err := ps.getGameIDAndColor(item[sortKey], colorAV)
			if err != nil {
				return model.Player{}, err
			}

			p.Games[gID] = color
			continue
		}

		name, ok := ps.getPlayerName(item[sortKey], item[ps.getNameKey()])
		if !ok {
			return model.Player{}, errors.New(`got unexpected payload`)
		} else if p.Name != `` {
			return model.Player{}, errors.New(`found two names`)
		}

		p.Name = name
	}

	if p.Name == `` {
		// This player _must_ have had at least a name stored, otherwise
		// we done messed up
		return model.Player{}, persistence.ErrPlayerNotFound
	}

	return p, nil
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
	data := map[string]types.AttributeValue{
		partitionKey: &types.AttributeValueMemberS{
			Value: string(p.ID),
		},
		sortKey: &types.AttributeValueMemberS{
			Value: getSortKeyPrefix(ps),
		},
		ps.getNameKey(): &types.AttributeValueMemberS{
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
		ConditionExpression: condExpr,
	})
	if err != nil {
		if isConditionalError(err) {
			return persistence.ErrPlayerAlreadyExists
		}
		return err
	}

	return nil
}

func (ps *playerService) BeginGame(gID model.GameID, players []model.Player) error {
	for _, p := range players {
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
	return playerColorAttributeName
}

func (ps *playerService) getNameKey() string {
	return playerNameAttributeName
}
