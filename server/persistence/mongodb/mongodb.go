package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb/mapbson"
)

const (
	dbName     string = `cribbage`
	gamesCol   string = `games`
	playersCol string = `players`
)

var _ persistence.DB = (*mongodb)(nil)

type mongodb struct {
	bsonRegistry *bsoncodec.Registry
	client       *mongo.Client
}

func New(uri string) (persistence.DB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &mongodb{
		client:       client,
		bsonRegistry: mapbson.CustomRegistry(),
	}, nil
}

func (m *mongodb) GetGame(id model.GameID) (model.Game, error) {
	result := model.Game{}
	gColOpts := []*options.CollectionOptions{
		{
			Registry: m.bsonRegistry,
		},
	}
	collection := m.client.Database(dbName).Collection(gamesCol, gColOpts...)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)

	if err != nil {
		return model.Game{}, err
	}

	return result, nil
}

func (m *mongodb) GetPlayer(id model.PlayerID) (model.Player, error) {
	result := model.Player{}
	gColOpts := []*options.CollectionOptions{
		{
			Registry: m.bsonRegistry,
		},
	}
	collection := m.client.Database(dbName).Collection(playersCol, gColOpts...)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)

	if err != nil {
		return model.Player{}, err
	}
	return result, nil
}

func (m *mongodb) GetInteraction(id model.PlayerID) (interaction.Player, error) {
	return nil, errors.New(`unimplemented`)
}

func (m *mongodb) SaveGame(g model.Game) error {
	gColOpts := []*options.CollectionOptions{
		{
			Registry: m.bsonRegistry,
		},
	}
	collection := m.client.Database(dbName).Collection(gamesCol, gColOpts...)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// TODO insert this like the memory does
	_, err := collection.InsertOne(ctx, g)
	return err
}

func (m *mongodb) CreatePlayer(p model.Player) error {
	gColOpts := []*options.CollectionOptions{
		{
			Registry: m.bsonRegistry,
		},
	}
	collection := m.client.Database(dbName).Collection(playersCol, gColOpts...)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// check if the player already exists
	sr := collection.FindOne(ctx, bson.M{"id": p.ID})
	if sr.Err() != mongo.ErrNoDocuments {
		return errors.New(`player already exists`)
	}

	_, err := collection.InsertOne(ctx, p)
	return err
}

func (m *mongodb) AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	return nil

}

func (m *mongodb) SaveInteraction(i interaction.Player) error {
	return nil
}
