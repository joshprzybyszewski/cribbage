//nolint:dupl
package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	playerCollectionIndex string = `id`
)

var _ persistence.PlayerService = (*playerService)(nil)

type playerService struct {
	ctx     context.Context
	session mongo.Session
	col     *mongo.Collection
}

func getPlayerService(
	ctx context.Context,
	session mongo.Session,
	mdb *mongo.Database,
	r *bsoncodec.Registry,
) (persistence.PlayerService, error) {

	col := mdb.Collection(playersCollectionName, &options.CollectionOptions{
		Registry: r,
	})

	idxs := col.Indexes()
	hasIndex, err := hasPlayerCollectionIndex(ctx, idxs)
	if err != nil {
		return nil, err
	}
	if !hasIndex {
		err = createPlayerCollectionIndex(ctx, idxs)
		if err != nil {
			return nil, err
		}
	}

	return &playerService{
		ctx:     ctx,
		session: session,
		col:     col,
	}, nil
}

func hasPlayerCollectionIndex(ctx context.Context, idxs mongo.IndexView) (bool, error) {
	return hasCollectionIndex(ctx, idxs, playerCollectionIndex)
}

func createPlayerCollectionIndex(ctx context.Context, idxs mongo.IndexView) error {
	return createCollectionIndex(ctx, idxs, playerCollectionIndex)
}

func bsonPlayerIDFilter(id model.PlayerID) interface{} {
	return bson.M{playerCollectionIndex: id} // model.Player.ID
}

func (ps *playerService) Get(id model.PlayerID) (model.Player, error) {
	result := model.Player{}
	filter := bsonPlayerIDFilter(id)
	err := mongo.WithSession(ps.ctx, ps.session, func(sc mongo.SessionContext) error {
		return ps.col.FindOne(sc, filter).Decode(&result)
	})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Player{}, persistence.ErrPlayerNotFound
		}
		return model.Player{}, err
	}
	return result, nil
}

func (ps *playerService) Create(p model.Player) error {
	// check if the player already exists
	filter := bsonPlayerIDFilter(p.ID)
	err := mongo.WithSession(ps.ctx, ps.session, func(sc mongo.SessionContext) error {
		c, err := ps.col.Find(sc, filter)
		if err != nil {
			return err
		}
		if c.Next(sc) {
			return persistence.ErrPlayerAlreadyExists
		}
		return nil
	})
	if err != nil {
		return err
	}

	return mongo.WithSession(ps.ctx, ps.session, func(sc mongo.SessionContext) error {
		_, err := ps.col.InsertOne(sc, p)
		// TODO could check the returned result to see how we did
		return err
	})
}

func (ps *playerService) UpdateGameColor(pID model.PlayerID, gID model.GameID, color model.PlayerColor) error {
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

	p.Games[gID] = color

	filter := bsonPlayerIDFilter(pID)
	opt := &options.FindOneAndReplaceOptions{}
	opt.SetUpsert(true)
	return mongo.WithSession(ps.ctx, ps.session, func(sc mongo.SessionContext) error {
		sr := ps.col.FindOneAndReplace(sc, filter, p)
		// TODO could check the returned result to see how we did
		return sr.Err()
	})
}
