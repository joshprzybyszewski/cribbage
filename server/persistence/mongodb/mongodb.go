package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DB = (*mongodb)(nil)

type mongodb struct {
	client *mongo.Client
}

func New() (persistence.DB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	return &mongodb{
		client: client,
	}, nil
}

func (m *mongodb) GetGame(id model.GameID) (model.Game, error) {
	return model.Game{}, errors.New(`does not have player`)
}

func (m *mongodb) GetPlayer(id model.PlayerID) (model.Player, error) {
	return model.Player{}, errors.New(`does not have player`)
}

func (m *mongodb) GetInteraction(id model.PlayerID) (interaction.Player, error) {
	return nil, errors.New(`does not have player`)
}

func (m *mongodb) SaveGame(g model.Game) error {
	return nil
}

func (m *mongodb) CreatePlayer(p model.Player) error {
	return nil
}

func (m *mongodb) AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	return nil

}

func (m *mongodb) SaveInteraction(i interaction.Player) error {
	return nil
}
