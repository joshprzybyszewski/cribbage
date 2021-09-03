//nolint:dupl
package dynamo

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	ctx context.Context

	svc *dynamodb.DynamoDB
}

func getInteractionService(
	ctx context.Context,
	svc *dynamodb.DynamoDB,
) (persistence.InteractionService, error) {

	return &interactionService{
		ctx: ctx,
		svc: svc,
	}, nil
}

func (s *interactionService) Get(id model.PlayerID) (interaction.PlayerMeans, error) {
	return interaction.PlayerMeans{}, errors.New(`todo`)
	/*
		svc := dynamodb.New(session.New())
		// I want to minimize the number of dynamo tables I use:
		// "You should maintain as few tables as possible in a DynamoDB application."
		// -https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-general-nosql-design.html
		input := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				dbName: {
					Keys: []map[string]*dynamodb.AttributeValue{{
						partitionKey: &dynamodb.AttributeValue{
							S: aws.String(string(id)),
						},
						sortKey: &dynamodb.AttributeValue{
							S: aws.String(string(dynamoInteractionServiceSortKey)),
						},
						// TODO use a "sort key" that defines this as the "interaction" model for the player
					}},
					// TODO figure out what the projexp should be for this?
					ProjectionExpression: aws.String("max(idk)"),
				},
			},
		}

		result2, err := svc.BatchGetItem(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result2)

		result := interaction.PlayerMeans{}
		filter := bsonInteractionFilter(id)
		err = mongo.WithSession(s.ctx, s.session, func(sc mongo.SessionContext) error {
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
	*/
}

func (s *interactionService) Create(pm interaction.PlayerMeans) error {
	return errors.New(`todo`)
	/*
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
	*/
}

func (s *interactionService) Update(pm interaction.PlayerMeans) error {
	return errors.New(`todo`)
	/*
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
	*/
}
