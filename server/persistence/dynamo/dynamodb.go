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
	svc := getDynamoService(ctx, df.endpoint)

	dto, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(dbName),
	})
	if err != nil || dto == nil {
		return nil, fmt.Errorf("DescribeTable ERROR: %v", err)
	}

	gs := newGameService(ctx, svc)
	ps := newPlayerService(ctx, svc)
	is := newInteractionService(ctx, svc)

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

func getDynamoService(
	ctx context.Context,
	endpoint string,
) *dynamodb.Client {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	cfg, err := config.LoadDefaultConfig(
		ctx,
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	opts := make([]func(o *dynamodb.Options), 0, 2)
	if len(endpoint) > 0 {
		// there should _only_ be an endpoint specified in local dev. Otherwise,
		// the magic aws config is supposed to figure it out.
		opts = append(opts,
			func(o *dynamodb.Options) {
				o.EndpointOptions = dynamodb.EndpointResolverOptions{
					DisableHTTPS: true,
				}
			},
			func(o *dynamodb.Options) {
				o.EndpointResolver = dynamodb.EndpointResolverFromURL(endpoint)
			},
		)
	}

	// Using the Config value, create the DynamoDB client
	return dynamodb.NewFromConfig(
		cfg,
		opts...,
	)
}

func (df dynamoFactory) Close() error {
	return nil
}

var _ persistence.DB = (*dynamoWrapper)(nil)

type dynamoWrapper struct {
	persistence.ServicesWrapper

	ctx context.Context
	// svc *dynamodb.Client // this may need to be used?
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
