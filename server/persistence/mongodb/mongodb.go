package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb/mapbson"
)

const (
	dbName                     string = `cribbage`
	gamesCollectionName        string = `games`
	playersCollectionName      string = `players`
	interactionsCollectionName string = `interactions`
)

func New(ctx context.Context, uri string) (persistence.DB, error) {
	if uri == `` {
		// If we don't know where to connect, use the default localhost URI
		uri = `mongodb://localhost:27017`
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	customRegistry := mapbson.CustomRegistry()
	mdb := client.Database(dbName)
	gs, err := getGameService(ctx, mdb, customRegistry)
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, mdb, customRegistry)
	if err != nil {
		return nil, err
	}
	is, err := getInteractionService(ctx, mdb, customRegistry)
	if err != nil {
		return nil, err
	}

	return persistence.New(
		gs,
		ps,
		is,
	), nil
}
