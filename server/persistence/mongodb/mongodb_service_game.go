//nolint:dupl
package mongodb

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	gameCollectionIndex string = `gameID`
)

type gameList struct {
	GameID model.GameID `bson:"gameID"`
	Games  []model.Game `bson:"games,omitempty"`
}

type persistedGameList struct {
	GameID    model.GameID `bson:"gameID"`
	TempGames []bson.M     `bson:"games,omitempty"`
}

func bsonGameIDFilter(id model.GameID) interface{} {
	// This filter should satisfy each of these fields being set:
	// gameList.GameID = id
	// persistedGameList.GameID = id
	return bson.M{gameCollectionIndex: id}
}

type getGameOptions struct {
	latest  bool
	all     bool
	actions map[int]struct{}
}

var _ persistence.GameService = (*gameService)(nil)

type gameService struct {
	ctx     context.Context
	session mongo.Session
	col     *mongo.Collection
}

func getGameService(
	ctx context.Context,
	session mongo.Session,
	mdb *mongo.Database,
	r *bsoncodec.Registry,
) (persistence.GameService, error) {

	col := mdb.Collection(gamesCollectionName, &options.CollectionOptions{
		Registry: r,
	})

	idxs := col.Indexes()
	hasIndex, err := hasGameCollectionIndex(ctx, idxs)
	if err != nil {
		return nil, err
	}
	if !hasIndex {
		err = createGameCollectionIndex(ctx, idxs)
		if err != nil {
			return nil, err
		}
	}

	return &gameService{
		ctx:     ctx,
		session: session,
		col:     col,
	}, nil
}

func hasGameCollectionIndex(ctx context.Context, idxs mongo.IndexView) (bool, error) {
	return hasCollectionIndex(ctx, idxs, gameCollectionIndex)
}

func createGameCollectionIndex(ctx context.Context, idxs mongo.IndexView) error {
	return createCollectionIndex(ctx, idxs, gameCollectionIndex)
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

func (gs *gameService) getSingleGame(id model.GameID, opts getGameOptions) (model.Game, error) {
	games, err := gs.getGameStates(id, opts)
	if err != nil {
		return model.Game{}, err
	}
	if len(games) != 1 {
		return model.Game{}, errors.New(`action doesn't exist`)
	}
	return games[0], nil
}

func (gs *gameService) getGameStates(id model.GameID, opts getGameOptions) ([]model.Game, error) {
	pgl := persistedGameList{}
	filter := bsonGameIDFilter(id)

	err := mongo.WithSession(gs.ctx, gs.session, func(sc mongo.SessionContext) error {
		err := gs.col.FindOne(sc, filter).Decode(&pgl)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return persistence.ErrGameNotFound
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if opts.actions == nil {
		opts.actions = make(map[int]struct{})
	}
	if opts.latest {
		opts.actions[len(pgl.TempGames)-1] = struct{}{}
	}

	gl := gameList{
		GameID: id,
		Games:  make([]model.Game, 0, len(pgl.TempGames)),
	}

	for i, tempGame := range pgl.TempGames {
		if _, ok := opts.actions[i]; !ok && !opts.all {
			continue
		}

		obj, err := json.Marshal(tempGame)
		if err != nil {
			return nil, err
		}

		g, err := jsonutils.UnmarshalGame(obj)
		if err != nil {
			return nil, err
		}

		gl.Games = append(gl.Games, g)
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

func (gs *gameService) Save(g model.Game) error {
	saved := gameList{}
	filter := bsonGameIDFilter(g.ID)

	err := mongo.WithSession(gs.ctx, gs.session, func(sc mongo.SessionContext) error {
		return gs.col.FindOne(sc, filter).Decode(&saved)
	})

	if err != nil {
		// if this is the first time saving the game, then we get ErrNoDocuments
		if err != mongo.ErrNoDocuments {
			return err
		}

		// Since this is the first save, we should have _no_ actions
		if len(g.Actions) != 0 {
			return persistence.ErrGameInitialSave
		}

		saved.GameID = g.ID
		saved.Games = []model.Game{g}

		return mongo.WithSession(gs.ctx, gs.session, func(sc mongo.SessionContext) error {
			var ior *mongo.InsertOneResult
			ior, err = gs.col.InsertOne(sc, saved)
			if err != nil {
				return err
			}
			if ior.InsertedID == nil {
				// not sure if this is the right thing to check
				return errors.New(`game not saved`)
			}

			return nil
		})
	}

	if saved.GameID != g.ID {
		return errors.New(`bad save somewhere`)
	}
	err = validateGameState(saved.Games, g)
	if err != nil {
		return err
	}

	saved.Games = append(saved.Games, g)

	return gs.saveGameList(saved)
}

func validateGameState(savedGames []model.Game, newGameState model.Game) error {
	if len(savedGames) != len(newGameState.Actions) {
		return persistence.ErrGameActionsOutOfOrder
	}
	for i := range savedGames {
		savedActions := savedGames[i].Actions
		myKnownActions := newGameState.Actions[:i]
		if len(savedActions) != len(myKnownActions) {
			return persistence.ErrGameActionsOutOfOrder
		}
		for ai, a := range savedActions {
			if a.ID != myKnownActions[ai].ID || a.Overcomes != myKnownActions[ai].Overcomes {
				return persistence.ErrGameActionsOutOfOrder
			}
		}
	}
	return nil
}

func (gs *gameService) saveGameList(saved gameList) error {
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
}
