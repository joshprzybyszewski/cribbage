package mongodb

import (
	"context"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// needs to match interaction.PlayerMeans.PlayerID
	interactionCollectionIndex string = `playerID`
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	ctx context.Context
	col *mongo.Collection
}

func getInteractionService(
	ctx context.Context,
	mdb *mongo.Database,
	r *bsoncodec.Registry,
) (persistence.InteractionService, error) {

	col := mdb.Collection(interactionsCollectionName, &options.CollectionOptions{
		Registry: r,
	})

	idxs := col.Indexes()
	hasIndex, err := hasInteractionCollectionIndex(ctx, idxs)
	if err != nil {
		return nil, err
	}
	if !hasIndex {
		err = createInteractionCollectionIndex(ctx, idxs)
		if err != nil {
			return nil, err
		}
	}

	return &interactionService{
		ctx: ctx,
		col: col,
	}, nil
}

func hasInteractionCollectionIndex(ctx context.Context, idxs mongo.IndexView) (bool, error) {
	return hasCollectionIndex(ctx, idxs, interactionCollectionIndex)
}

func createInteractionCollectionIndex(ctx context.Context, idxs mongo.IndexView) error {
	return createCollectionIndex(ctx, idxs, interactionCollectionIndex)
}

func bsonInteractionFilter(id model.PlayerID) interface{} {
	// interaction.PlayerMeans{PlayerID: id}
	return bson.M{`playerID`: id}
}

func (s *interactionService) Get(id model.PlayerID) (interaction.PlayerMeans, error) {
	result := interaction.PlayerMeans{}
	filter := bsonInteractionFilter(id)
	err := s.col.FindOne(s.ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return interaction.PlayerMeans{}, persistence.ErrInteractionNotFound
		}
		return interaction.PlayerMeans{}, err
	}

	return result, nil
}

func (s *interactionService) Create(pm interaction.PlayerMeans) error {
	_, err := s.Get(pm.PlayerID)
	if err != nil && err != persistence.ErrInteractionNotFound {
		return err
	}

	_, err = s.col.InsertOne(s.ctx, pm)
	return err
}

func (s *interactionService) Update(pm interaction.PlayerMeans) error {
	if _, err := s.Get(pm.PlayerID); err == persistence.ErrInteractionNotFound {
		_, err = s.col.InsertOne(s.ctx, pm)
		return err
	}

	opt := &options.ReplaceOptions{}
	opt.SetUpsert(true)

	_, err := s.col.ReplaceOne(s.ctx, pm, opt)
	return err
}
