package dynamo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

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

var _ persistence.DBFactory = dynamoFactory{}

type dynamoFactory struct {
	uri string
}

func NewFactory(uri string) (persistence.DBFactory, error) {
	svc := dynamodb.New(session.New())
	return dynamoFactory{
		uri: uri,
	}, nil
}

func (df dynamoFactory) New(ctx context.Context) (persistence.DB, error) {
	uri := df.uri
	if uri == `` {
		// The default URI without replicas used to be:
		// `mongodb://localhost:27017`
		// But now we should be running with replicaset, so let's talk to all three
		// Followed instructions here: http://thecodebarbarian.com/introducing-run-rs-zero-config-mongodb-runner
		uri = `mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs`
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Always start a session so we can create a transaction if needed
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

	dw := dynamoWrapper{
		ServicesWrapper: sw,
		ctx:             ctx,
		client:          client,
		session:         sess,
	}

	return &dw, nil
}

func (df dynamoFactory) Close() error {
	return nil
}

var _ persistence.DB = (*dynamoWrapper)(nil)

type dynamoWrapper struct {
	persistence.ServicesWrapper

	ctx     context.Context
	client  *mongo.Client
	session mongo.Session
}

func (dw *dynamoWrapper) Close() error {
	return dw.client.Disconnect(dw.ctx)
}

func (dw *dynamoWrapper) Start() error {
	if dw.session == nil {
		return errors.New(`no session to use`)
	}

	txOpts := options.Transaction()
	txOpts.SetReadConcern(readconcern.Local())
	txOpts.SetReadPreference(readpref.Primary())
	txOpts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	mct := maxCommitTime
	txOpts.SetMaxCommitTime(&mct)

	return dw.session.StartTransaction(txOpts)
}

func (dw *dynamoWrapper) Commit() error {
	return dw.finishTx(func(sc mongo.SessionContext) error {
		return dw.session.CommitTransaction(sc)
	})
}

func (dw *dynamoWrapper) Rollback() error {
	return dw.finishTx(func(sc mongo.SessionContext) error {
		return dw.session.AbortTransaction(sc)
	})
}

func (dw *dynamoWrapper) finishTx(finisher func(mongo.SessionContext) error) (err error) {
	if dw.session == nil {
		return errors.New(`missing session`)
	}

	defer func() {
		if err == nil {
			// only end session if there was no error
			dw.session.EndSession(dw.ctx)
		}
	}()

	err = mongo.WithSession(dw.ctx, dw.session, finisher)
	return err
}
