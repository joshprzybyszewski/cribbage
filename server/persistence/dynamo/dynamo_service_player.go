//nolint:dupl
package dynamo

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

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

	svc *dynamodb.DynamoDB
}

func getPlayerService(
	ctx context.Context,
	svc *dynamodb.DynamoDB,
) (persistence.PlayerService, error) {

	return &playerService{
		ctx: ctx,
		svc: svc,
	}, nil
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			dbName: {
				Keys: []map[string]*dynamodb.AttributeValue{{
					partitionKey: &dynamodb.AttributeValue{
						S: aws.String(string(id)),
					},
					// I want the sort key to be prefix-able
					sortKey: &dynamodb.AttributeValue{
						S: aws.String(string(playerServiceSortKeyPrefix)),
					},
					// TODO figure out what the getter should be to only nab the relevant info for the player?
				}},
				// TODO figure out what the projexp should be?
				ProjectionExpression: aws.String("max(idk)"),
			},
		},
	}

	dynamoResult, err := ps.svc.BatchGetItem(input)
	if err != nil {
		fmt.Println(err)
		return model.Player{}, err
	}
	fmt.Println(dynamoResult)
	return model.Player{}, errors.New(`josh TODO`)
}

func (ps *playerService) Create(p model.Player) error {
	// TODO create the player in the put item input
	input := &dynamodb.PutItemInput{}

	output, err := ps.svc.PutItem(input)
	if err != nil {
		return err
	}
	if output.Attributes != nil {
		// TODO validate output?
		return persistence.ErrPlayerAlreadyExists
	}
	return nil
}

func (ps *playerService) BeginGame(gID model.GameID, players []model.Player) error {
	// TODO for each player in players, create an item that includes the game ID
	return nil
}

func (ps *playerService) UpdateGameColor(pID model.PlayerID, gID model.GameID, color model.PlayerColor) error {
	// TODO, get the pID gID combo and assign the color
	p, err := ps.Get(pID)
	if err != nil {
		return err
	}

	if p.Games == nil {
		p.Games = make(map[model.GameID]model.PlayerColor, 1)
	}

	if c, ok := p.Games[gID]; ok {
		if c != color {
			return errors.New(`mismatched player-games color`)
		}

		// Nothing to do; the player already knows its color
		return nil
	}

	// TODO create the playergame in the put item input
	// p.Games[gID] = color
	input := &dynamodb.PutItemInput{}

	output, err := ps.svc.PutItem(input)
	if err != nil {
		return err
	}
	if output.Attributes != nil {
		// TODO validate output?
		return persistence.ErrPlayerAlreadyExists
	}
	return nil
}
