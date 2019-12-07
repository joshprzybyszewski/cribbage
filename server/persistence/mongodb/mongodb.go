package mongodb

import (
	"context"
	"encoding/json"
	"errors"
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

const (
	latestGameAction int = -1
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

	// needs to match savedGame.ID
	err = m.setupIndex(ctx, `gameID`, m.gamesCollection().Indexes())
	if err != nil {
		return nil, err
	}

	// needs to match model.Player.ID
	err = m.setupIndex(ctx, `id`, m.playersCollection().Indexes())
	if err != nil {
		return nil, err
	}

	// needs to match interaction.PlayerMeans.PlayerID
	err = m.setupIndex(ctx, `playerID`, m.interactionsCollection().Indexes())
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

type gameList struct {
	GameID model.GameID `bson:"gameID"`
	Games  []model.Game `bson:"games,omitempty"`
}

type persistedGameList struct {
	GameID    model.GameID `bson:"gameID"`
	TempGames []bson.M     `bson:"games,omitempty"`
}

func (m *mongodb) GetGame(id model.GameID) (model.Game, error) {
	return m.getGameAtAction(id, latestGameAction)
}

func (m *mongodb) GetGameAction(id model.GameID, numActions uint) (model.Game, error) {
	return m.getGameAtAction(id, int(numActions))
}

func (m *mongodb) getGameAtAction(id model.GameID, numActions int) (model.Game, error) {
	pgl := persistedGameList{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := persistedGameList{GameID: id}
	err := m.gamesCollection().FindOne(ctx, filter).Decode(&pgl)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Game{}, persistence.ErrGameNotFound
		}
		return model.Game{}, err
	}

	if numActions == latestGameAction {
		numActions = len(pgl.TempGames) - 1
	}

	if numActions < 0 || numActions >= len(pgl.TempGames) {
		return model.Game{}, errors.New(`action doesn't exist`)
	}

	tempGame := pgl.TempGames[numActions]
	obj, err := json.Marshal(tempGame)
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
		if err == mongo.ErrNoDocuments {
			return model.Player{}, persistence.ErrPlayerNotFound
		}
		return model.Player{}, err
	}
	return result, nil
}

func (m *mongodb) SaveGame(g model.Game) error {
	// TODO make this transactional
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	saved := gameList{}
	filter := gameList{GameID: g.ID}
	err := m.gamesCollection().FindOne(ctx, filter).Decode(&saved)
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
		_, err = m.gamesCollection().InsertOne(ctx, saved)
		return err
	}

	if saved.GameID != g.ID {
		return errors.New(`bad save somewhere`)
	}
	if len(saved.Games) != len(g.Actions) {
		// TODO we could do a deeper check on the actions
		// i.e. saved.games == g.Action[:len(g.Actions)-1]
		return persistence.ErrGameActionsOutOfOrder
	}

	saved.Games = append(saved.Games, g)

	_, err = m.gamesCollection().ReplaceOne(ctx, filter, saved)
	return err
}

func (m *mongodb) CreatePlayer(p model.Player) error {
	collection := m.playersCollection()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// check if the player already exists
	filter := bson.M{`id`: p.ID} // model.Player{ID: p.ID}
	c, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	if c.Next(ctx) {
		return persistence.ErrPlayerAlreadyExists
	}

	_, err = collection.InsertOne(ctx, p)
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

	if _, err := m.GetInteraction(i.PlayerID); err == persistence.ErrInteractionNotFound {
		_, err = collection.InsertOne(ctx, i)
		return err
	}

	opt := &options.ReplaceOptions{}
	opt.SetUpsert(true)

	_, err := collection.ReplaceOne(ctx, i, opt)
	return err
}

func (m *mongodb) GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error) {
	pm := interaction.PlayerMeans{}
	filter := bson.M{`playerID`: id} // interaction.PlayerMeans{PlayerID: id}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := m.interactionsCollection().FindOne(ctx, filter).Decode(&pm)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return interaction.PlayerMeans{}, persistence.ErrInteractionNotFound
		}
		return interaction.PlayerMeans{}, err
	}
	return pm, nil
}
