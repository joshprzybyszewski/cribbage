package mapbson

import (
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"

	"github.com/joshprzybyszewski/cribbage/model"
)

func CustomRegistry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()

	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	registerPlayerGames(rb)
	registerGameBlockingPlayers(rb)
	registerGamePlayerColors(rb)
	registerGameHands(rb)
	registerGameScores(rb)

	return rb.Build()
}

func registerPlayerGames(rb *bsoncodec.RegistryBuilder) {
	var p model.Player
	var gamesType reflect.Type = reflect.TypeOf(p.Games)

	c := newCustomMapCoder(gamesType, gameIDStringEncoder, gameIDStringDecoder)

	rb.RegisterEncoder(gamesType, c)
	rb.RegisterDecoder(gamesType, c)
}

func gameIDStringEncoder(k reflect.Value) (string, error) {
	gID := k.Interface().(model.GameID)
	return strconv.Itoa(int(gID)), nil
}

func gameIDStringDecoder(key string) (reflect.Value, error) {
	i, err := strconv.Atoi(key)
	if err != nil {
		return reflect.Value{}, err
	}
	typedKey := model.GameID(i)
	return reflect.ValueOf(typedKey), nil
}

func registerGameBlockingPlayers(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var blockingPlayersType reflect.Type = reflect.TypeOf(g.BlockingPlayers)

	c := newCustomMapCoder(blockingPlayersType, playerIDStringEncoder, playerIDStringDecoder)

	rb.RegisterEncoder(blockingPlayersType, c)
	rb.RegisterDecoder(blockingPlayersType, c)
}

func registerGamePlayerColors(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var playerColorsType reflect.Type = reflect.TypeOf(g.PlayerColors)

	c := newCustomMapCoder(playerColorsType, playerIDStringEncoder, playerIDStringDecoder)

	rb.RegisterEncoder(playerColorsType, c)
	rb.RegisterDecoder(playerColorsType, c)
}

func registerGameHands(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var handsType reflect.Type = reflect.TypeOf(g.Hands)

	c := newCustomMapCoder(handsType, playerIDStringEncoder, playerIDStringDecoder)

	rb.RegisterEncoder(handsType, c)
	rb.RegisterDecoder(handsType, c)
}

func registerGameScores(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var gameScoresType reflect.Type = reflect.TypeOf(g.CurrentScores)

	c := newCustomMapCoder(gameScoresType, playerColorStringEncoder, playerColorStringDecoder)

	rb.RegisterEncoder(gameScoresType, c)
	rb.RegisterDecoder(gameScoresType, c)
}

func playerIDStringEncoder(k reflect.Value) (string, error) {
	pID := k.Interface().(model.PlayerID)
	return string(pID), nil
}

func playerIDStringDecoder(key string) (reflect.Value, error) {
	typedKey := model.PlayerID(key)
	return reflect.ValueOf(typedKey), nil
}

func playerColorStringEncoder(k reflect.Value) (string, error) {
	pc := k.Interface().(model.PlayerColor)
	return strconv.Itoa(int(pc)), nil
}

func playerColorStringDecoder(key string) (reflect.Value, error) {
	n, err := strconv.Atoi(key)
	if err != nil {
		return reflect.Value{}, err
	}
	typedKey := model.PlayerColor(n)
	return reflect.ValueOf(typedKey), nil
}
