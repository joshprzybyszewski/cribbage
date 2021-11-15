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
	gameBytesAttributeName = `gameBytes`
)

var _ persistence.GameService = (*gameService)(nil)

type gameService struct {
	ctx context.Context

	svc *dynamodb.Client
}

func newGameService(
	ctx context.Context,
	svc *dynamodb.Client,
) persistence.GameService {
	return &gameService{
		ctx: ctx,
		svc: svc,
	}
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

type getGameOptions struct {
	latest      bool
	actionIndex uint
}

func (gs *gameService) getGame(
	id model.GameID,
	opts getGameOptions,
) (model.Game, error) {
	pkName := `:gID`
	pk := strconv.Itoa(int(id))
	skName := `:sk`
	sk := gs.getSpecForAllGameActions()
	if !opts.latest {
		sk = gs.getSpecForGameActionIndex(opts.actionIndex)
	}
	hp := hasPrefix{
		pkName: pkName,
		skName: skName,
	}

	sif := false

	qi := &dynamodb.QueryInput{
		ScanIndexForward:       &sif,
		TableName:              aws.String(dbName),
		KeyConditionExpression: hp.conditionExpression(),
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
		return model.Game{}, err
	}
	if len(qo.Items) == 0 {
		return model.Game{}, persistence.ErrGameNotFound
	}

	item := qo.Items[0]
	if !opts.latest {
		// make sure that the index we got back matches the one we requested
		spec, ok := item[sortKey].(*types.AttributeValueMemberS)
		if !ok {
			return model.Game{}, persistence.ErrGameActionDecode
		}
		i, err := gs.getGameActionIndexFromSpec(spec.Value)
		if err != nil {
			return model.Game{}, err
		}
		if i != int(opts.actionIndex) {
			return model.Game{}, errors.New(`retrieved unexpected game action index`)
		}
	}

	gb, ok := item[gs.getSerGameKey()].(*types.AttributeValueMemberB)
	if !ok {
		return model.Game{}, persistence.ErrGameActionDecode
	}
	return jsonutils.UnmarshalGame(gb.Value)
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

	return gs.writeGame(writeGameOptions{
		game:        g,
		actionIndex: uint(len(g.Actions)),
		overwrite:   true,
	})
}

func (gs *gameService) Begin(g model.Game) error {
	return gs.writeGame(writeGameOptions{
		game:        g,
		actionIndex: 0,
	})
}

func (gs *gameService) Save(g model.Game) error {
	err := persistence.ValidateLatestActionBelongs(g)
	if err != nil {
		return err
	}

	// validate that the actions on this game are known by the previous game.
	sg, err := gs.getGame(g.ID, getGameOptions{
		latest: true,
	})
	if err != nil {
		return err
	}
	if len(sg.Actions)+1 != len(g.Actions) {
		// The new game state can only have one additional action
		return persistence.ErrGameActionsOutOfOrder
	}
	for i := range sg.Actions {
		if !actionsAreEqual(sg.Actions[i], g.Actions[i]) {
			return persistence.ErrGameActionsOutOfOrder
		}
	}

	return gs.writeGame(writeGameOptions{
		game:        g,
		actionIndex: uint(len(sg.Actions) + 1),
	})
}

func actionsAreEqual(a, b model.PlayerAction) bool {
	return a.GameID == b.GameID &&
		a.ID == b.ID &&
		a.Overcomes == b.Overcomes &&
		a.TimestampStr == b.TimestampStr
}

type writeGameOptions struct {
	game        model.Game
	actionIndex uint
	overwrite   bool
}

// writeGame will write the given game and action
// This method assumes you've already done game state validation.
func (gs *gameService) writeGame(opts writeGameOptions) error {
	obj, err := json.Marshal(opts.game)
	if err != nil {
		return err
	}

	pii := &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{
				Value: strconv.Itoa(int(opts.game.ID)),
			},
			sortKey: &types.AttributeValueMemberS{
				Value: gs.getSpecForGameActionIndex(opts.actionIndex),
			},
			gs.getSerGameKey(): &types.AttributeValueMemberB{
				Value: obj,
			},
		},
	}

	if opts.overwrite {
		// we want to find out if we overwrote items, so specify ReturnValues
		pii.ReturnValues = types.ReturnValueAllOld
	} else {
		// Use a conditional expression to only write items if this
		// <HASH:RANGE> tuple doesn't already exist.
		// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
		// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
		pii.ConditionExpression = notExists{}.conditionExpression()
	}

	pio, err := gs.svc.PutItem(gs.ctx, pii)
	if err != nil {
		if isConditionalError(err) {
			return persistence.ErrGameActionsOutOfOrder
		}
		return err
	}

	if opts.overwrite {
		// We need to check that we actually overwrote an element
		if _, ok := pio.Attributes[gs.getSerGameKey()]; !ok {
			// oh no! We wanted to overwrite a game, but we didn't!
			return persistence.ErrGameActionsOutOfOrder
		}
	}

	return nil
}

func (gs *gameService) getSerGameKey() string {
	return gameBytesAttributeName
}

func (gs *gameService) getSpecForAllGameActions() string {
	return getSortKeyPrefix(gs) + `@`
}

func (gs *gameService) getSpecForGameActionIndex(i uint) string {
	// Since we print out leading zeros to six places, we could
	// have issues if our cribbage games ever take more than
	// 999,999 actions. This is an arbitrary limit and one that
	// we're unlikely to ever encounter.
	return gs.getSpecForAllGameActions() + fmt.Sprintf(`%06d`, i)
}

func (gs *gameService) getGameActionIndexFromSpec(s string) (int, error) {
	s = strings.TrimPrefix(s, gs.getSpecForAllGameActions())
	return strconv.Atoi(s)
}
