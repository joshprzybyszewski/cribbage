package mongodb

import (
	"context"
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
	dbName          string = `cribbage`
	gamesCol        string = `games`
	playersCol      string = `players`
	interactionsCol string = `interactions`
)

var _ persistence.DB = (*mongodb)(nil)

type mongodb struct {
	bsonRegistry *bsoncodec.Registry
	client       *mongo.Client
}

func New(uri string) (persistence.DB, error) {
	if uri == `` {
		// If we don't know where to connect, use the default localhost URI
		uri = `mongodb://localhost:27017`
	}
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

func (m *mongodb) gamesCollection() *mongo.Collection {
	gColOpts := []*options.CollectionOptions{{
		Registry: m.bsonRegistry,
	}}
	return m.client.Database(dbName).Collection(gamesCol, gColOpts...)
}

func (m *mongodb) playersCollection() *mongo.Collection {
	gColOpts := []*options.CollectionOptions{{
		Registry: m.bsonRegistry,
	}}
	return m.client.Database(dbName).Collection(playersCol, gColOpts...)
}

func (m *mongodb) interactionsCollection() *mongo.Collection {
	gColOpts := []*options.CollectionOptions{{
		Registry: m.bsonRegistry,
	}}
	return m.client.Database(dbName).Collection(interactionsCol, gColOpts...)
}

func (m *mongodb) GetGame(id model.GameID) (model.Game, error) {
	result := model.Game{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := m.gamesCollection().FindOne(ctx, bson.M{"id": id}).Decode(&result)

	if err != nil {
		return model.Game{}, err
	}

	return result, nil
}

func (m *mongodb) GetPlayer(id model.PlayerID) (model.Player, error) {
	result := model.Player{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := m.playersCollection().FindOne(ctx, bson.M{"id": id}).Decode(&result)

	if err != nil {
		return model.Player{}, err
	}
	return result, nil
}

func (m *mongodb) SaveGame(g model.Game) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// TODO insert this like the memory does, where we append it to the end of a list of games,
	// in order of how many actions they've had
	_, err := m.gamesCollection().InsertOne(ctx, g)
	return err
}

func (m *mongodb) CreatePlayer(p model.Player) error {
	collection := m.playersCollection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// check if the player already exists
	sr := collection.FindOne(ctx, bson.M{"id": p.ID})
	if sr.Err() != mongo.ErrNoDocuments {
		return persistence.ErrPlayerAlreadyExists
	}

	_, err := collection.InsertOne(ctx, p)
	return err
}

func (m *mongodb) AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	collection := m.playersCollection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// Overwrite the player's Games field with a new map
	// TODO figure out replacement
	var replacement interface{}
	collection.FindOneAndReplace(ctx, bson.M{"id": id}, replacement)
	// TODO do this stufffr
	return nil

}

func (m *mongodb) SaveInteraction(i interaction.PlayerMeans) error {
	collection := m.interactionsCollection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// check if the player already exists
	sr := collection.FindOne(ctx, bson.M{"id": i.PlayerID})
	if sr.Err() != mongo.ErrNoDocuments {
		return persistence.ErrPlayerAlreadyExists
	}

	_, err := collection.InsertOne(ctx, i)
	return err
}

func (m *mongodb) GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error) {
	pm := interaction.PlayerMeans{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := m.interactionsCollection().FindOne(ctx, bson.M{"id": id}).Decode(&pm)

	if err != nil {
		return interaction.PlayerMeans{}, err
	}
	return pm, nil
}
