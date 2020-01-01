package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb/mapbson"
)

const (
	dbName                     string = `cribbage`
	gamesCollectionName        string = `games`
	playersCollectionName      string = `players`
	interactionsCollectionName string = `interactions`
)

const (
	maxCommitTime time.Duration = 10 * time.Second // something very large for now -- this should be reduced
)

var _ persistence.DB = (*mongoWrapper)(nil)

type mongoWrapper struct {
	persistence.ServicesWrapper

	ctx     context.Context
	client  *mongo.Client
	session mongo.Session
}

func New(ctx context.Context, uri string) (persistence.DB, error) {
	if uri == `` {
		// If we don't know where to connect, use the default localhost URI
		// uri = `mongodb://localhost:27017`
		// Now we should be running with replicaset, so let's talk to all three
		// Followed instructions here: http://thecodebarbarian.com/introducing-run-rs-zero-config-mongodb-runner
		uri = `mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs`
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// start a session so we can create a transaction
	sess, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	customRegistry := mapbson.CustomRegistry()
	mdb := client.Database(dbName)
	gs, err := getGameService(ctx, sess, mdb, customRegistry)
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, sess, mdb, customRegistry)
	if err != nil {
		return nil, err
	}
	is, err := getInteractionService(ctx, sess, mdb, customRegistry)
	if err != nil {
		return nil, err
	}

	sw := persistence.NewServicesWrapper(
		gs,
		ps,
		is,
	)

	mw := mongoWrapper{
		ServicesWrapper: sw,
		ctx:             ctx,
		client:          client,
		session:         sess,
	}

	return &mw, nil
}

func (mw *mongoWrapper) Start() error {
	if mw.session == nil {
		return errors.New(`no session to use`)
	}

	txOpts := &options.TransactionOptions{}
	txOpts.SetReadConcern(readconcern.Local())                         // local
	txOpts.SetReadPreference(readpref.Primary())                       // primary
	txOpts.SetWriteConcern(writeconcern.New(writeconcern.WMajority())) // majority
	mct := maxCommitTime
	txOpts.SetMaxCommitTime(&mct)

	return mw.session.StartTransaction(txOpts)
}

func (mw *mongoWrapper) Commit() (err error) {
	return mw.finishTx(func(sc mongo.SessionContext) error {
		return mw.session.CommitTransaction(sc)
	})
}

func (mw *mongoWrapper) Rollback() (err error) {
	return mw.finishTx(func(sc mongo.SessionContext) error {
		return mw.session.AbortTransaction(sc)
	})
}

func (mw *mongoWrapper) finishTx(finisher func(mongo.SessionContext) error) (err error) {
	if mw.session == nil {
		return errors.New(`missing session`)
	}

	defer func() {
		if err != nil {
			// only end session & disconnect client if there was no error
			return
		}
		mw.session.EndSession(mw.ctx)
		if err = mw.client.Disconnect(mw.ctx); err != nil {
			fmt.Printf("got error on disconnect: %+v\n", err)
		}
	}()

	err = mongo.WithSession(mw.ctx, mw.session, finisher)
	return err
}
