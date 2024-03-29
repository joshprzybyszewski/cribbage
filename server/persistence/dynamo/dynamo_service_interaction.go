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

const (
	interactionInfoAttributeName   = `info`
	interactionPreferAttributeName = `prefer`
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	ctx context.Context

	svc *dynamodb.Client
}

func newInteractionService(
	ctx context.Context,
	svc *dynamodb.Client,
) persistence.InteractionService {
	return &interactionService{
		ctx: ctx,
		svc: svc,
	}
}

func (is *interactionService) Get(
	id model.PlayerID,
) (interaction.PlayerMeans, error) {
	pkName := `:ipID`
	pk := string(id)
	skName := `:sk`
	sk := getSortKeyPrefix(is)
	hp := hasPrefix{
		pkName: pkName,
		skName: skName,
	}

	createQuery := newQueryInputFactory(getQueryInputParams(
		pk, pkName, sk, skName, hp.conditionExpression(),
	))
	items, err := fullQuery(is.ctx, is.svc, createQuery)
	if err != nil {
		return interaction.PlayerMeans{}, err
	}

	return is.buildPlayerMeansFromItems(id, items)
}

func (is *interactionService) buildPlayerMeansFromItems(
	id model.PlayerID,
	items []map[string]types.AttributeValue,
) (interaction.PlayerMeans, error) {

	pm := interaction.PlayerMeans{
		PlayerID:      id,
		PreferredMode: interaction.Unknown,
	}

	for _, item := range items {
		if preferAV, ok := item[is.getPreferKey()]; ok {
			if pm.PreferredMode != interaction.Unknown {
				return interaction.PlayerMeans{}, errors.New(`preferred mode already set`)
			}

			preferAVN, ok := preferAV.(*types.AttributeValueMemberN)
			if !ok {
				return interaction.PlayerMeans{}, errors.New(`wrong prefer attribute type`)
			}
			preferMode, err := strconv.Atoi(preferAVN.Value)
			if err != nil {
				return interaction.PlayerMeans{}, err
			}

			pm.PreferredMode = interaction.Mode(preferMode)
			continue
		}

		mode, err := is.getInteractionMode(
			item[sortKey],
		)
		if err != nil {
			return interaction.PlayerMeans{}, err
		}

		m := interaction.Means{
			Mode: mode,
		}

		serInfo, err := is.getSerInfo(
			item[is.getInfoKey()],
		)
		if err != nil {
			return interaction.PlayerMeans{}, err
		}
		err = m.AddSerializedInfo(serInfo)
		if err != nil {
			return interaction.PlayerMeans{}, err
		}

		pm.Interactions = append(pm.Interactions, m)

	}

	return pm, nil
}

func (is *interactionService) getInteractionMode(
	specAV types.AttributeValue,
) (interaction.Mode, error) {
	specAVS, ok := specAV.(*types.AttributeValueMemberS)
	if !ok {
		return interaction.Unknown, errors.New(`wrong spec`)
	}

	mode, err := is.getInteractionMeansModeFromSpec(specAVS.Value)
	if err != nil {
		return interaction.Unknown, err
	}
	return mode, nil
}

func (is *interactionService) getSerInfo(
	infoAV types.AttributeValue,
) ([]byte, error) {
	infoAVB, ok := infoAV.(*types.AttributeValueMemberB)
	if !ok {
		return nil, errors.New(`wrong info type`)
	}

	return infoAVB.Value, nil
}

func (is *interactionService) Create(pm interaction.PlayerMeans) error {
	return is.write(writePlayerMeansOptions{
		pm:         pm,
		isCreation: true,
	})
}

func (is *interactionService) Update(pm interaction.PlayerMeans) error {
	return is.write(writePlayerMeansOptions{
		pm: pm,
	})
}

type writePlayerMeansOptions struct {
	pm         interaction.PlayerMeans
	isCreation bool
}

func (is *interactionService) write(opts writePlayerMeansOptions) error {
	data := map[string]types.AttributeValue{
		partitionKey: &types.AttributeValueMemberS{
			Value: string(opts.pm.PlayerID),
		},
		sortKey: &types.AttributeValueMemberS{
			Value: getSortKeyPrefix(is),
		},
		is.getPreferKey(): &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(opts.pm.PreferredMode)),
		},
	}

	pii := &dynamodb.PutItemInput{
		TableName: aws.String(dbName),
		Item:      data,
	}

	if opts.isCreation {
		// Use a conditional expression to only write items if this
		// <HASH:RANGE> tuple doesn't already exist.
		// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
		// and https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
		pii.ConditionExpression = notExists{}.conditionExpression()
	}

	_, err := is.svc.PutItem(is.ctx, pii)
	if err != nil {
		if isConditionalError(err) {
			return persistence.ErrInteractionAlreadyExists
		}
		return err
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
	s = strings.TrimPrefix(s, getSortKeyPrefix(is)+`|`)
	i, err := strconv.Atoi(s)
	if err != nil {
		return interaction.Unknown, err
	}
	return interaction.Mode(i), nil
}

func (is *interactionService) getInfoKey() string {
	return interactionInfoAttributeName
}

func (is *interactionService) getPreferKey() string {
	return interactionPreferAttributeName
}
