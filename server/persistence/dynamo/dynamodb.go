package dynamo

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	dbName       = `cribbage`
	partitionKey = `DDBid`
	sortKey      = `spec`
)

var _ persistence.DBFactory = dynamoFactory{}

type dynamoFactory struct {
	endpoint string
}

func NewFactory(endpoint string) (persistence.DBFactory, error) {
	return dynamoFactory{
		endpoint: endpoint,
	}, nil
}

func (df dynamoFactory) New(ctx context.Context) (persistence.DB, error) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	cfg, err := config.LoadDefaultConfig(
		ctx,
		// TODO how should I re-set region?
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
			o.EndpointResolver = dynamodb.EndpointResolverFromURL(df.endpoint)
		},
	)

	dto, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(dbName),
	})
	if err != nil || dto == nil {
		return nil, fmt.Errorf("DescribeTable ERROR: %v", err)
	}

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
