package dynamo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	dbName       = `cribbage`
	partitionKey = `DDBid`
	sortKey      = `spec`

	dynamoPlayerServiceSortKey      = playerServiceSortKeyPrefix
	dynamoInteractionServiceSortKey = `interaction`

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
	return dynamoFactory{
		uri: uri,
	}, nil
}

func (df dynamoFactory) New(ctx context.Context) (persistence.DB, error) {
	// TODO figure out this junk
	os.Setenv(`AWS_ACCESS_KEY_ID`, `DUMMYIDEXAMPLE`)
	os.Setenv(`AWS_SECRET_ACCESS_KEY`, `DUMMYEXAMPLEKEY`)
	endpoint := `http://localhost:18079`

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(
		cfg,
		func(o *dynamodb.Options) {
			o.EndpointOptions = dynamodb.EndpointResolverOptions{
				DisableHTTPS: true, // todo remove this when we deploy non-locally
			}
		},
		func(o *dynamodb.Options) {
			// TODO don't do this in non-local
			o.EndpointResolver = dynamodb.EndpointResolverFromURL(endpoint)
		},
	)

	tn := dbName
	dto, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &tn,
	})
	fmt.Printf("DescribeTable output, err : %#v, %+v\n", dto.Table, err)

	gs, err := getGameService(ctx, svc)
	if err != nil {
		return nil, err
	}
	ps, err := getPlayerService(ctx, svc)
	if err != nil {
		return nil, err
	}
	is, err := getInteractionService(ctx, svc)
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
	}

	return &dw, nil
}

func (df dynamoFactory) Close() error {
	return nil
}

var _ persistence.DB = (*dynamoWrapper)(nil)

type dynamoWrapper struct {
	persistence.ServicesWrapper

	ctx context.Context
	svc *dynamodb.Client
}

func (dw *dynamoWrapper) Close() error {
	// TODO I don't think there's anything to do?
	return nil
}

func (dw *dynamoWrapper) Start() error {
	// TODO figure out transactionality in dynamodb
	return nil
}

func (dw *dynamoWrapper) Commit() error {
	// TODO figure out transactionality in dynamodb
	return nil
}

func (dw *dynamoWrapper) Rollback() error {
	// TODO figure out transactionality in dynamodb
	return nil
}
