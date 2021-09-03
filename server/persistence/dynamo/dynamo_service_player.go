//nolint:dupl
package dynamo

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

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
				// ProjectionExpression: aws.String("max(idk)"),
			},
		},
	}

	dynamoResult, err := ps.svc.BatchGetItem(input)
	if err != nil {
		return model.Player{}, err
	}
	fmt.Println(dynamoResult)
	for i, resp := range dynamoResult.Responses[dbName] {
		dp := dynamoPlayer{}
		// TODO unmarshalling a map doesn't work (i.e. the games)
		dynamodbattribute.UnmarshalMap(resp, &dp)
		fmt.Printf("\ti, resp := dp\n\t%d, %+v := %+v\n\t%+v\n", i, resp, dp, dp.Player)
	}
	return model.Player{}, errors.New(`josh TODO`)
}

type dynamoPlayer struct {
	ID   string `json:"DDBid"`
	Spec string `json:"spec"`

	Player model.Player `json:"serPlayer"`
}

func (ps *playerService) Create(p model.Player) error {
	if len(p.Games) > 0 {
		return errors.New(`you cannot create a player that is _already_ in games!`)
	}

	data, err := dynamodbattribute.MarshalMap(dynamoPlayer{
		ID:     string(p.ID),
		Spec:   playerServiceSortKeyPrefix, // TODO this is going to have a game id at the end for player colors!
		Player: p,
	})
	if err != nil {
		return err
	}

	output, err := ps.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item:      data,
	})
	if err != nil {
		return err
	}
	// TODO find a way to discover if the player already existed.
	if output.Attributes != nil {
		// TODO validate output?
		return persistence.ErrPlayerAlreadyExists
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
		data, err := dynamodbattribute.MarshalMap(dynamoPlayerInGame{
			ID:    string(p.ID),
			Spec:  fmt.Sprintf("%s%d", playerServiceGameSortKey, gID),
			Color: model.UnsetColor,
		})
		if err != nil {
			return err
		}

		// TODO do these in separate goroutines (aka parallelize)
		output, err := ps.svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(dbName),
			Item:      data,
		})
		if err != nil {
			return err
		}
		// TODO find a way to discover if the <playerID:gameID> already existed.
		if output.Attributes != nil {
			// TODO validate output?
			return persistence.ErrPlayerAlreadyExists
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

	data, err := dynamodbattribute.MarshalMap(dynamoPlayerInGame{
		ID:    string(p.ID),
		Spec:  fmt.Sprintf("%s%d", playerServiceGameSortKey, gID),
		Color: color,
	})
	if err != nil {
		return err
	}

	output, err := ps.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item:      data,
	})
	if err != nil {
		return err
	}
	if output.Attributes != nil {
		// TODO validate output?
		return persistence.ErrPlayerAlreadyExists
	}
	return nil
}
