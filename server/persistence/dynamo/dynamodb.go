package dynamo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	os.Setenv(`AWS_ACCESS_KEY_ID`, `DUMMYIDEXAMPLE`)
	os.Setenv(`AWS_SECRET_ACCESS_KEY`, `DUMMYEXAMPLEKEY`)
	endpoint := `http://localhost:18079`
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: &endpoint,
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tn := dbName
	dto, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &tn,
	})
	fmt.Printf("DescribeTable output, err : %+v, %+v\n", dto, err)

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
	svc *dynamodb.DynamoDB
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
