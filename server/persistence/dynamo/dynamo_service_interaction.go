//nolint:dupl
package package dynamo

import (
	"context"
	"errors"

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
	ctx     context.Context
	session mongo.Session
	col     *mongo.Collection
}

func getInteractionService(
	ctx context.Context,
	session mongo.Session,
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
		ctx:     ctx,
		session: session,
		col:     col,
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
	svc := dynamodb.New(session.New())
	// I want to minimize the number of dynamo tables I use:
	// "You should maintain as few tables as possible in a DynamoDB application."
	// -https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-general-nosql-design.html
	dynamoGamesTableName := `cribbage`
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			dynamoGamesTableName: {
				Keys: []map[string]*dynamodb.AttributeValue{{
					"PlayerID": &dynamodb.AttributeValue{
						S: aws.String(string(id)),
					},
					// TODO use a "sort key" that defines this as the "interaction" model for the player
				}},
				// TODO figure out what the projexp should be for this?
				ProjectionExpression: aws.String("max(idk)"),
			},
		},
	}

	result, err := svc.BatchGetItem(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	result := interaction.PlayerMeans{}
	filter := bsonInteractionFilter(id)
	err := mongo.WithSession(s.ctx, s.session, func(sc mongo.SessionContext) error {
		err := s.col.FindOne(sc, filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return persistence.ErrInteractionNotFound
			}
			return err
		}
		return nil
	})
	if err != nil {
		return interaction.PlayerMeans{}, err
	}

	return result, nil
}

func (s *interactionService) Create(pm interaction.PlayerMeans) error {
	_, err := s.Get(pm.PlayerID)
	if err != nil && err != persistence.ErrInteractionNotFound {
		return err
	}

	return mongo.WithSession(s.ctx, s.session, func(sc mongo.SessionContext) error {
		ior, err := s.col.InsertOne(sc, pm)
		if err != nil {
			return err
		}
		if ior.InsertedID == nil {
			// :shrug: not sure if this is the right thing to check
			return errors.New(`interaction not created`)
		}

		return nil
	})
}

func (s *interactionService) Update(pm interaction.PlayerMeans) error {
	if _, err := s.Get(pm.PlayerID); err == persistence.ErrInteractionNotFound {
		return mongo.WithSession(s.ctx, s.session, func(sc mongo.SessionContext) error {
			ior, err := s.col.InsertOne(sc, pm)
			if err != nil {
				return err
			}
			if ior.InsertedID == nil {
				// :shrug: not sure if this is the right thing to check
				return errors.New(`interaction not updated`)
			}

			return nil
		})
	}

	opt := &options.ReplaceOptions{}
	opt.SetUpsert(true)

	return mongo.WithSession(s.ctx, s.session, func(sc mongo.SessionContext) error {
		ur, err := s.col.ReplaceOne(sc, pm, opt)
		if err != nil {
			return err
		}

		switch {
		case ur.ModifiedCount > 1:
			return errors.New(`modified too many interactions`)
		case ur.MatchedCount > 1:
			return errors.New(`matched more than one interaction`)
		case ur.UpsertedCount > 1:
			return errors.New(`upserted more than one interaction`)
		}

		return nil
	})
}
