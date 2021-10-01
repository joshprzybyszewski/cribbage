//nolint:dupl
package dynamo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	gameServiceSortKeyPrefix string = `game`
)

type gameList struct {
	GameID model.GameID `bson:"gameID"`
	Games  []model.Game `bson:"games,omitempty"`
}

func (gl *gameList) add(i int, g model.Game) {
	// TODO
	for i >= len(gl.Games) {
		gl.Games = append(gl.Games, model.Game{})
	}
	gl.Games[i] = g
}

type getGameOptions struct {
	latest  bool
	all     bool
	actions map[int]struct{}
}

var _ persistence.GameService = (*gameService)(nil)

type gameService struct {
	ctx context.Context

	svc *dynamodb.Client
}

func getGameService(
	ctx context.Context,
	svc *dynamodb.Client,
) (persistence.GameService, error) {

	return &gameService{
		ctx: ctx,
		svc: svc,
	}, nil
}

func (gs *gameService) Get(id model.GameID) (model.Game, error) {
	return gs.getSingleGame(id, getGameOptions{
		latest: true,
	})
}

func (gs *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	return gs.getSingleGame(id, getGameOptions{
		actions: map[int]struct{}{int(numActions): {}},
	})
}

func (gs *gameService) getSingleGame(
	id model.GameID,
	opts getGameOptions,
) (model.Game, error) {

	games, err := gs.getGameStates(id, opts)
	if err != nil {
		return model.Game{}, err
	}
	if len(games) != 1 {
		return model.Game{}, errors.New(`action doesn't exist`)
	}
	return games[0], nil
}

func (gs *gameService) getGameStates(id model.GameID, opts getGameOptions) ([]model.Game, error) { // nolint:gocyclo
	gl := gameList{}

	// I want to minimize the number of dynamo tables I use:
	// "You should maintain as few tables as possible in a DynamoDB application."
	// -https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-general-nosql-design.html
	// input := &dynamodb.BatchGetItemInput{
	// 	RequestItems: map[string]*dynamodb.KeysAndAttributes{
	// 		dbName: {
	// 			Keys: []map[string]*dynamodb.AttributeValue{{
	// 				partitionKey: &dynamodb.AttributeValue{
	// 					S: aws.String(string(id)),
	// 				},
	// 				// TODO I don't remember right now how to get based on partition/sort key
	// 			}},
	// 			// TODO if only latest, then we can use a projexp to filter down
	// 			ProjectionExpression: aws.String("max(numGameActions)"),
	// 		},
	// 	},
	// }

	// result, err := gs.svc.BatchGetItem(gs.ctx, input)
	// if err != nil {

	tableName := dbName
	pkName := `:gID`
	pk := strconv.Itoa(int(id))
	skName := `:sk`
	sk := gameServiceSortKeyPrefix
	keyCondExpr := fmt.Sprintf("DDBid = %s and begins_with(spec, %s)", pkName, skName)
	qo, err := gs.svc.Query(gs.ctx, &dynamodb.QueryInput{
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
		fmt.Println(err)
	}
	fmt.Printf("query for games: %+v\n", qo)

	if opts.actions == nil {
		opts.actions = make(map[int]struct{})
	}

	var highestI int
	var highestG model.Game

	for _, item := range qo.Items {
		spec, ok := item[`spec`].(*types.AttributeValueMemberS)
		if !ok {
			return nil, persistence.ErrGameActionDecode
		}
		i, err := getGameActionIndexFromSpec(spec.Value)
		if err != nil {
			return nil, err
		}

		gb, ok := item[`gameBytes`].(*types.AttributeValueMemberB)
		if !ok {
			return nil, persistence.ErrGameActionDecode
		}
		g, err := jsonutils.UnmarshalGame(gb.Value)
		if err != nil {
			return nil, err
		}
		if i > highestI {
			highestI = i
			highestG = g
		}

		if _, ok := opts.actions[i]; !ok && !opts.all {
			continue
		}

		gl.add(i, g)
	}

	if opts.latest {
		gl.add(highestI, highestG)
	}

	return gl.Games, nil
}

func (gs *gameService) UpdatePlayerColor(gID model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	g, err := gs.Get(gID)
	if err != nil {
		return err
	}

	if c, ok := g.PlayerColors[pID]; ok {
		if c != color {
			return errors.New(`mismatched game-player color`)
		}

		// the Game already knows this player's color; nothing to do
		return nil
	}

	games, err := gs.getGameStates(gID, getGameOptions{
		all: true,
	})
	if err != nil {
		return err
	}

	recentGame := games[len(games)-1]
	if recentGame.PlayerColors == nil {
		recentGame.PlayerColors = make(map[model.PlayerID]model.PlayerColor, 1)
	}
	recentGame.PlayerColors[pID] = color

	games[len(games)-1] = recentGame
	newGameList := gameList{
		GameID: gID,
		Games:  games,
	}

	return gs.saveGameList(newGameList)
}

func (gs *gameService) Begin(g model.Game) error {
	return gs.Save(g)
}

func (gs *gameService) Save(g model.Game) error {

	games, err := gs.getGameStates(g.ID, getGameOptions{
		all: true,
	})
	if err != nil {
		return err
	}

	saved := gameList{
		GameID: g.ID,
		Games:  games,
	}
	if saved.GameID != g.ID {
		return errors.New(`bad save somewhere`)
	}

	// TODO does this still belong?
	// err = persistence.ValidateLatestActionBelongs(g)
	// if err != nil {
	// 	return err
	// }

	saved.Games = append(saved.Games, g)

	return gs.saveGameList(saved)
}

func getSpecForGameActionIndex(i int) string {
	return gameServiceSortKeyPrefix + `@` + strconv.Itoa(i)
}

func getGameActionIndexFromSpec(s string) (int, error) {
	s = strings.TrimPrefix(s, gameServiceSortKeyPrefix+`@`)
	return strconv.Atoi(s)
}

func (gs *gameService) saveGameList(saved gameList) error {
	// craft an id that uses len(saved.Games)

	obj, err := json.Marshal(saved.Games[len(saved.Games)-1])
	if err != nil {
		return err
	}

	data := map[string]types.AttributeValue{
		`DDBid`: &types.AttributeValueMemberS{
			Value: strconv.Itoa(int(saved.GameID)),
		},
		`spec`: &types.AttributeValueMemberS{
			Value: gameServiceSortKeyPrefix + `@` + strconv.Itoa(len(saved.Games)),
		},
		`gameBytes`: &types.AttributeValueMemberB{
			Value: obj,
		},
	}

	// Use a conditional expression to only write items if this
	// <HASH:RANGE> tuple doesn't already exist.
	// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
	// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
	condExpr := `attribute_not_exists(DDBid) AND attribute_not_exists(spec)`

	fmt.Printf("playerService.Create data = %+v\n", data)

	_, err = gs.svc.PutItem(gs.ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(dbName),
		Item:                data,
		ConditionExpression: &condExpr,
	})
	if err != nil {
		switch err.(type) {
		case *types.ConditionalCheckFailedException:
			return persistence.ErrGameActionSave
		}
		return err
	}

	return nil
	/*
		filter := bsonGameIDFilter(saved.GameID)
		return mongo.WithSession(gs.ctx, gs.session, func(sc mongo.SessionContext) error {
			ur, err := gs.col.ReplaceOne(sc, filter, saved)
			if err != nil {
				return err
			}

			switch {
			case ur.ModifiedCount > 1:
				return errors.New(`modified too many games`)
			case ur.MatchedCount > 1:
				return errors.New(`matched more than one game entry`)
			case ur.UpsertedCount > 1:
				return errors.New(`replaced more than one game`)
			}

			return nil
		})
	*/
}
