//nolint:dupl
package dynamo

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	ctx context.Context

	svc *dynamodb.Client
}

func getInteractionService(
	ctx context.Context,
	svc *dynamodb.Client,
) (persistence.InteractionService, error) {

	return &interactionService{
		ctx: ctx,
		svc: svc,
	}, nil
}

func (is *interactionService) Get(
	id model.PlayerID,
) (interaction.PlayerMeans, error) {
	ret := interaction.PlayerMeans{
		PlayerID:      id,
		PreferredMode: interaction.Unknown,
	}

	pkName := `:pID`
	pk := string(id)
	skName := `:sk`
	sk := getSortKeyPrefix(is)
	keyCondExpr := getConditionExpression(equalsID, pkName, hasPrefix, skName)

	createQuery := func() *dynamodb.QueryInput {
		return &dynamodb.QueryInput{
			TableName:              aws.String(dbName),
			KeyConditionExpression: &keyCondExpr,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				pkName: &types.AttributeValueMemberS{
					Value: pk,
				},
				skName: &types.AttributeValueMemberS{
					Value: sk,
				},
			},
		}
	}
	qi := createQuery()

	for {
		qo, err := is.svc.Query(is.ctx, qi)
		if err != nil {
			return interaction.PlayerMeans{}, err
		}

		err = is.populatePlayerMeansFromItems(&ret, qo.Items)
		// check LastEvaluatedKey to know if we need to paginate the responses
		// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.Pagination
		if len(qo.LastEvaluatedKey) == 0 {
			break
		}

		qi = createQuery()
		qi.ExclusiveStartKey = qo.LastEvaluatedKey
	}

	return ret, nil
}

func (is *interactionService) populatePlayerMeansFromItems(
	pm *interaction.PlayerMeans,
	items []map[string]types.AttributeValue,
) error {
	for _, item := range items {

		if preferAV, ok := item[`prefer`]; ok {
			preferAVN, ok := preferAV.(*types.AttributeValueMemberN)
			if !ok {
				return errors.New(`wrong prefer`)
			}
			preferMode, err := strconv.Atoi(preferAVN.Value)
			if err != nil {
				return err
			}
			pm.PreferredMode = interaction.Mode(preferMode)
			continue
		}

		mode, serInfo, err := is.getInteractionModeAndSerInfo(item[sortKey], item[is.getInfoKey()])
		if err != nil {
			// invalid persisted means
			return err
		}

		m := interaction.Means{
			Mode: mode,
		}
		m.AddSerializedInfo(serInfo)

		pm.Interactions = append(pm.Interactions, m)

	}

	return nil
}

func (is *interactionService) getInteractionModeAndSerInfo(
	specAV, infoAV types.AttributeValue,
) (interaction.Mode, []byte, error) {
	specAVS, ok := specAV.(*types.AttributeValueMemberS)
	if !ok {
		return interaction.Unknown, nil, errors.New(`wrong spec`)
	}

	mode, err := is.getInteractionMeansModeFromSpec(specAVS.Value)
	if err != nil {
		return interaction.Unknown, nil, err
	}

	infoAVN, ok := infoAV.(*types.AttributeValueMemberB)
	if !ok {
		return interaction.Unknown, nil, errors.New(`wrong info type`)
	}

	return mode, infoAVN.Value, nil
}

func (is *interactionService) Create(pm interaction.PlayerMeans) error {
	return is.write(writePlayerMeansOptions{
		pm: pm,
	})
}

func (is *interactionService) Update(pm interaction.PlayerMeans) error {
	return is.write(writePlayerMeansOptions{
		pm:        pm,
		overwrite: true,
	})
}

type writePlayerMeansOptions struct {
	pm        interaction.PlayerMeans
	overwrite bool
}

func (is *interactionService) write(opts writePlayerMeansOptions) error {
	data := map[string]types.AttributeValue{
		partitionKey: &types.AttributeValueMemberS{
			Value: string(opts.pm.PlayerID),
		},
		sortKey: &types.AttributeValueMemberS{
			Value: getSortKeyPrefix(is),
		},
		`prefer`: &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(opts.pm.PreferredMode)),
		},
	}

	pii := &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item:      data,
	}

	if opts.overwrite {
		// we want to find out if we overwrote items, so specify ReturnValues
		pii.ReturnValues = types.ReturnValueAllOld
	} else {
		// Use a conditional expression to only write items if this
		// <HASH:RANGE> tuple doesn't already exist.
		// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
		// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
		condExpr := getConditionExpression(notExists, ``, notExists, ``)
		pii.ConditionExpression = &condExpr
	}

	pio, err := is.svc.PutItem(is.ctx, pii)
	if err != nil {
		switch err.(type) {
		case *types.ConditionalCheckFailedException:
			return persistence.ErrInteractionAlreadyExists
		}
		return err
	}

	if opts.overwrite {
		// We need to check that we actually overwrote an element
		if _, ok := pio.Attributes[`spec`]; !ok {
			// oh no! We wanted to overwrite a game, but we didn't!
			return persistence.ErrInteractionUnexpected
		}
	}

	for _, m := range opts.pm.Interactions {
		pii, err = is.getPutItemInputForMeans(opts.pm.PlayerID, m)
		if err != nil {
			return err
		}
		_, err = is.svc.PutItem(is.ctx, pii)
		if err != nil {
			return err
		}
	}

	return nil
}

func (is *interactionService) getPutItemInputForMeans(
	playerID model.PlayerID,
	m interaction.Means,
) (*dynamodb.PutItemInput, error) {
	info, err := m.GetSerializedInfo()
	if err != nil {
		return nil, err
	}

	return &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item: map[string]types.AttributeValue{
			partitionKey: &types.AttributeValueMemberS{
				Value: string(playerID),
			},
			sortKey: &types.AttributeValueMemberS{
				Value: is.getSpecForInteractionMeans(m),
			},
			is.getInfoKey(): &types.AttributeValueMemberB{
				Value: info,
			},
		},
	}, nil
}

func (is *interactionService) getSpecForInteractionMeans(
	m interaction.Means,
) string {
	return getSortKeyPrefix(is) + `|` + strconv.Itoa(int(m.Mode))
}

func (is *interactionService) getInteractionMeansModeFromSpec(s string) (interaction.Mode, error) {
	s = strings.TrimPrefix(s, getSortKeyPrefix(is)+`@`)
	i, err := strconv.Atoi(s)
	if err != nil {
		return interaction.Unknown, err
	}
	return interaction.Mode(i), nil
}

func (is *interactionService) getInfoKey() string {
	return `info`
}
