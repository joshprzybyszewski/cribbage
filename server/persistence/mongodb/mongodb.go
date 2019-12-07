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

	// needs to match gameList.GameID
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
	gs, err := m.getGameStates(id, map[int]struct{}{numActions: {}})
	if err != nil {
		return model.Game{}, err
	}
	if len(gs) != 1 {
		return model.Game{}, errors.New(`action doesn't exist`)
	}
	return gs[0], nil
}

func (m *mongodb) getGameStates(id model.GameID, actionStates map[int]struct{}) ([]model.Game, error) {
	pgl := persistedGameList{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{`gameID`: id} // persistedGameList{GameID: id}
	err := m.gamesCollection().FindOne(ctx, filter).Decode(&pgl)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, persistence.ErrGameNotFound
		}
		return nil, err
	}

	if _, ok := actionStates[latestGameAction]; ok {
		delete(actionStates, latestGameAction)
		actionStates[len(pgl.TempGames)-1] = struct{}{}
	}

	gl := gameList{
		GameID: id,
		Games:  make([]model.Game, 0, len(pgl.TempGames)),
	}

	for i, tempGame := range pgl.TempGames {
		if _, ok := actionStates[i]; actionStates != nil && !ok {
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
	filter := bson.M{`gameID`: g.ID} // gameList{GameID: g.ID}
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

	return m.saveGameList(ctx, saved)
}

func (m *mongodb) saveGameList(ctx context.Context, saved gameList) error {
	filter := bson.M{`gameID`: saved.GameID} // gameList{GameID: saved.GameID}
	_, err := m.gamesCollection().ReplaceOne(ctx, filter, saved)
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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := m.addPlayerColorToPlayer(ctx, id, color, gID)
	if err != nil {
		return err
	}

	err = m.addPlayerColorToGame(ctx, id, color, gID)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongodb) addPlayerColorToPlayer(
	ctx context.Context,
	id model.PlayerID,
	color model.PlayerColor,
	gID model.GameID,
) error {

	// Overwrite the player's Games field with a new map
	filter := bson.M{`id`: id} // model.Player{ID: p.ID}
	p, err := m.GetPlayer(id)
	if err != nil {
		return nil
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

	opt := &options.FindOneAndReplaceOptions{}
	opt.SetUpsert(true)
	sr := m.playersCollection().FindOneAndReplace(ctx, filter, p)
	return sr.Err()
}

func (m *mongodb) addPlayerColorToGame(
	ctx context.Context,
	id model.PlayerID,
	color model.PlayerColor,
	gID model.GameID,
) error {

	g, err := m.GetGame(gID)
	if err != nil {
		return nil
	}

	if c, ok := g.PlayerColors[id]; ok {
		if c != color {
			return errors.New(`mismatched game-player color`)
		}

		// the Game already knows this player's color; nothing to do
		return nil
	}

	games, err := m.getGameStates(gID, nil)
	if err != nil {
		return err
	}

	recentGame := games[len(games)-1]
	if recentGame.PlayerColors == nil {
		recentGame.PlayerColors = make(map[model.PlayerID]model.PlayerColor, 1)
	}
	recentGame.PlayerColors[id] = color

	games[len(games)-1] = recentGame
	newGameList := gameList{
		GameID: gID,
		Games:  games,
	}

	return m.saveGameList(ctx, newGameList)
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
