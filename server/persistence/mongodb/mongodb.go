package mongodb

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
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

	m := mongodb{
		client:       client,
		bsonRegistry: mapbson.CustomRegistry(),
	}

	err = m.setupIndex(ctx, `id`, m.gamesCollection().Indexes())
	if err != nil {
		return nil, err
	}

	err = m.setupIndex(ctx, `id`, m.playersCollection().Indexes())
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *mongodb) setupIndex(ctx context.Context, indexKey string, idxs mongo.IndexView) error {
	listMaxTimeOpt := &options.ListIndexesOptions{}
	listMaxTimeOpt.SetMaxTime(30 * time.Second)
	cur, err := idxs.List(ctx, listMaxTimeOpt)
	if err != nil {
		return err
	}

	for cur.Next(context.Background()) {
		index := bson.D{}
		err = cur.Decode(&index)
		if err != nil {
			return err
		}
		for _, i := range index {
			if key := i.Key; key == `key` {
				if val, ok := i.Value.(bson.D); ok && val[0].Key == indexKey {
					// found the desired index, exit
					return nil
				}
			}
		}
	}

	keys := bsonx.Doc{{Key: indexKey, Value: bsonx.Int64(int64(1))}}
	im := mongo.IndexModel{}
	im.Keys = keys
	opts := options.CreateIndexes().SetMaxTime(5 * time.Minute)

	_, err = idxs.CreateOne(ctx, im, opts)
	return err
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
	// var monthStatus interface{}
	// filter := bson.M{"_id": id}
	// tempResult := bson.M{}
	// err := db.Collection("Months").FindOne(ctx, filter).Decode(&tempResult)
	// if err == nil {
	// 	obj, _ := json.Marshal(tempResult)
	// 	err = json.Unmarshal(obj, &monthStatus)
	// }
	// return monthStatus, err

	// result := model.Game{}
	tempResult := bson.M{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := m.gamesCollection().FindOne(ctx, bson.M{"id": id}).Decode(&tempResult)

	if err != nil {
		return model.Game{}, err
	}

	obj, err := json.Marshal(tempResult)
	if err != nil {
		return model.Game{}, err
	}

	g, err := jsonutils.UnmarshalGame(obj)
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
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
	// TODO make this transactional
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	current := model.Game{}
	err := m.gamesCollection().FindOne(ctx, bson.M{"id": g.ID}).Decode(&current)
	if err != nil {
		// if this is the first time saving the game, then we get ErrNoDocuments
		if err != mongo.ErrNoDocuments {
			return err
		}

		// Since this is the first save, we should have _no_ actions
		if len(g.Actions) != 0 {
			return persistence.ErrGameInitialSave
		}

		_, err = m.gamesCollection().InsertOne(ctx, g)
		return err
	}

	if len(current.Actions)+1 != len(g.Actions) {
		// TODO we could do a deeper check on the actions
		// i.e. current.Actions == g.Action[:len(g.Actions)-1]
		return persistence.ErrGameActionsOutOfOrder
	}

	_, err = m.gamesCollection().ReplaceOne(ctx, bson.M{"id": g.ID}, g)
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
