//nolint:dupl
package dynamo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

type gameAtAction struct {
	Game        model.Game
	ActionIndex int
	Overwrite   bool
}

type getGameOptions struct {
	latest      bool
	actionIndex uint
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
	return gs.getGame(id, getGameOptions{
		latest: true,
	})
}

func (gs *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	return gs.getGame(id, getGameOptions{
		actionIndex: numActions,
	})
}

func (gs *gameService) getGame(
	id model.GameID,
	opts getGameOptions,
) (model.Game, error) { // nolint:gocyclo
	pkName := `:gID`
	pk := strconv.Itoa(int(id))
	skName := `:sk`
	sk := gameServiceSortKeyPrefix + `@`
	if !opts.latest {
		sk = getSpecForGameActionIndex(int(opts.actionIndex))
	}
	keyCondExpr := fmt.Sprintf("DDBid = %s and begins_with(spec, %s)", pkName, skName)

	// sif := false

	qi := &dynamodb.QueryInput{
		// ScanIndexForward:       &sif,
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
	qo, err := gs.svc.Query(gs.ctx, qi)
	if err != nil {
		fmt.Println(err)
		return model.Game{}, err
	}
	if len(qo.Items) == 0 {
		return model.Game{}, errors.New(`unexpected number of items returned`)
	}
	log.Printf("qo.LastEvaluatedKey: %+v\n", qo.LastEvaluatedKey)

	for i, item := range qo.Items {
		s := ``
		if spec, ok := item[`spec`].(*types.AttributeValueMemberS); ok {
			s = spec.Value
		}
		log.Printf("items[%d] = %T{%q}\n", i, item[`spec`], s)
	}

	item := qo.Items[0]
	if !opts.latest {
		// make sure that the index we got back matches the one we requested
		spec, ok := item[`spec`].(*types.AttributeValueMemberS)
		if !ok {
			return model.Game{}, persistence.ErrGameActionDecode
		}
		i, err := getGameActionIndexFromSpec(spec.Value)
		if err != nil {
			return model.Game{}, err
		}
		if i != int(opts.actionIndex) {
			return model.Game{}, errors.New(`retrieved unexpected game action index`)
		}
	}

	gb, ok := item[`gameBytes`].(*types.AttributeValueMemberB)
	if !ok {
		return model.Game{}, persistence.ErrGameActionDecode
	}
	g, err := jsonutils.UnmarshalGame(gb.Value)
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
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

	if g.PlayerColors == nil {
		g.PlayerColors = make(map[model.PlayerID]model.PlayerColor, 1)
	}
	g.PlayerColors[pID] = color

	return gs.saveGame(gameAtAction{
		Game:        g,
		ActionIndex: len(g.Actions) - 1,
		Overwrite:   true,
	})
}

func (gs *gameService) Begin(g model.Game) error {
	return gs.save(saveOptions{
		game:       g,
		isCreation: true,
	})
}

func (gs *gameService) Save(g model.Game) error {
	err := persistence.ValidateLatestActionBelongs(g)
	if err != nil {
		return err
	}
	return gs.save(saveOptions{
		game: g,
	})
}

type saveOptions struct {
	game       model.Game
	isCreation bool
}

func (gs *gameService) save(
	opts saveOptions,
) error {
	ai := 0
	if !opts.isCreation {
		sg, err := gs.getGame(opts.game.ID, getGameOptions{
			latest: true,
		})
		if err != nil {
			return err
		}
		ai = len(sg.Actions)
	}

	return gs.saveGame(gameAtAction{
		Game:        opts.game,
		ActionIndex: ai,
	})
}

func getSpecForGameActionIndex(i int) string {
	// Since we print out leading zeros to nine places, we could
	// have issues if our cribbage games ever take more than
	// 999,999,999 actions
	return gameServiceSortKeyPrefix + `@` + fmt.Sprintf("%09d", i)
}

func getGameActionIndexFromSpec(s string) (int, error) {
	s = strings.TrimPrefix(s, gameServiceSortKeyPrefix+`@`)
	return strconv.Atoi(s)
}

func (gs *gameService) saveGame(gaa gameAtAction) error {
	obj, err := json.Marshal(gaa.Game)
	if err != nil {
		return err
	}

	data := map[string]types.AttributeValue{
		`DDBid`: &types.AttributeValueMemberS{
			Value: strconv.Itoa(int(gaa.Game.ID)),
		},
		`spec`: &types.AttributeValueMemberS{
			Value: getSpecForGameActionIndex(gaa.ActionIndex),
		},
		`gameBytes`: &types.AttributeValueMemberB{
			Value: obj,
		},
	}

	fmt.Printf("gameService.saveGame data = %+v\n", data)
	fmt.Printf("gameService.saveGame gaa.Game = %#v\n", gaa.Game)

	pii := &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item:      data,
	}

	if gaa.Overwrite {
		// we want to find out if we overwrote items, so specify ReturnValues
		pii.ReturnValues = types.ReturnValueAllOld
	} else {
		// Use a conditional expression to only write items if this
		// <HASH:RANGE> tuple doesn't already exist.
		// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
		// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
		condExpr := `attribute_not_exists(DDBid) AND attribute_not_exists(spec)`

		pii.ConditionExpression = &condExpr
	}

	pio, err := gs.svc.PutItem(gs.ctx, pii)
	if err != nil {
		switch err.(type) {
		case *types.ConditionalCheckFailedException:
			return persistence.ErrGameActionSave
		}
		return err
	}

	if gaa.Overwrite {
		// We need to check that we actually overwrote an element
		if _, ok := pio.Attributes[`gameBytes`]; !ok {
			// oh no! We wanted to overwrite a game, but we didn't!
			return persistence.ErrGameActionSave
		}
	}

	return nil
}
