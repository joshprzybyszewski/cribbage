package mongodb

import (
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/joshprzybyszewski/cribbage/model"
)

func modelBSONRegistry() *bsoncodec.Registry {
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
	var gamesKeyType reflect.Type = reflect.TypeOf(model.InvalidGameID)
	var gamesValueType reflect.Type = reflect.TypeOf(model.Blue)

	c := &customMapCoder{
		customMapKind: gamesType.Kind(),
		keyKind:       gamesKeyType.Kind(),
		keyEncoder:    gameIDStringEncoder,
		keyDecoder:    gameIDStringDecoder,
		valueKind:     gamesValueType.Kind(),
	}

	rb.RegisterEncoder(gamesType, c)
	rb.RegisterDecoder(gamesType, c)
}

func registerGameBlockingPlayers(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var blockingPlayersType reflect.Type = reflect.TypeOf(g.BlockingPlayers)

	c := &customMapCoder{
		customMapKind: blockingPlayersType.Kind(),
		keyKind:       reflect.TypeOf(model.InvalidPlayerID).Kind(),
		valueKind:     reflect.TypeOf(model.DealCards).Kind(),
		keyEncoder:    playerIDStringEncoder,
		keyDecoder:    playerIDStringDecoder,
	}

	rb.RegisterEncoder(blockingPlayersType, c)
	rb.RegisterDecoder(blockingPlayersType, c)
}

func registerGamePlayerColors(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var playerColorsType reflect.Type = reflect.TypeOf(g.PlayerColors)

	c := &customMapCoder{
		customMapKind: playerColorsType.Kind(),
		keyKind:       reflect.TypeOf(model.InvalidPlayerID).Kind(),
		valueKind:     reflect.TypeOf(model.Blue).Kind(),
		keyEncoder:    playerIDStringEncoder,
		keyDecoder:    playerIDStringDecoder,
	}

	rb.RegisterEncoder(playerColorsType, c)
	rb.RegisterDecoder(playerColorsType, c)
}

func registerGameHands(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var handsType reflect.Type = reflect.TypeOf(g.Hands)
	var hand []model.Card

	c := &customMapCoder{
		customMapKind: handsType.Kind(),
		keyKind:       reflect.TypeOf(model.InvalidPlayerID).Kind(),
		valueKind:     reflect.TypeOf(hand).Kind(),
		keyEncoder:    playerIDStringEncoder,
		keyDecoder:    playerIDStringDecoder,
	}

	rb.RegisterEncoder(handsType, c)
	rb.RegisterDecoder(handsType, c)
}

func registerGameScores(rb *bsoncodec.RegistryBuilder) {
	var g model.Game
	var gameScoresType reflect.Type = reflect.TypeOf(g.CurrentScores)
	var score int

	c := &customMapCoder{
		customMapKind: gameScoresType.Kind(),
		keyKind:       reflect.TypeOf(model.Blue).Kind(),
		valueKind:     reflect.TypeOf(score).Kind(),
		keyEncoder:    playerColorStringEncoder,
		keyDecoder:    playerColorStringDecoder,
	}

	rb.RegisterEncoder(gameScoresType, c)
	rb.RegisterDecoder(gameScoresType, c)
}

var _ bsoncodec.ValueEncoder = (*customMapCoder)(nil)
var _ bsoncodec.ValueDecoder = (*customMapCoder)(nil)

type customMapCoder struct {
	keyKind   reflect.Kind
	valueKind reflect.Kind

	keyEncoder func(reflect.Value) (string, error)
	keyDecoder func(string) (reflect.Value, error)

	customMapKind reflect.Kind // should be reflect.TypeOf(model.Player{}.Games).Kind()
}

func (c *customMapCoder) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() ||
		val.Kind() != reflect.Map ||
		val.Type().Key().Kind() != c.keyKind ||
		val.Type().Elem().Kind() != c.valueKind {
		return bsoncodec.ValueEncoderError{Name: "PlayerGamesEncodeValue", Kinds: []reflect.Kind{c.customMapKind}, Received: val}
	}

	if val.IsNil() {
		err := vw.WriteNull()
		if err == nil {
			return nil
		}
	}

	dw, err := vw.WriteDocument()
	if err != nil {
		return err
	}

	encoder, err := ectx.LookupEncoder(val.Type().Elem())
	if err != nil {
		return err
	}

	keys := val.MapKeys()
	for _, k := range keys {
		keyString, err := c.keyEncoder(k)
		if err != nil {
			return err
		}

		vw, err := dw.WriteDocumentElement(keyString)
		if err != nil {
			return err
		}

		err = encoder.EncodeValue(ectx, vw, val.MapIndex(k))
		if err != nil {
			return err
		}
	}

	return dw.WriteDocumentEnd()
}

func (c *customMapCoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() ||
		val.Kind() != reflect.Map ||
		val.Type().Key().Kind() != c.keyKind ||
		val.Type().Elem().Kind() != c.valueKind {
		return bsoncodec.ValueDecoderError{Name: "PlayerGamesDecodeValue", Kinds: []reflect.Kind{c.customMapKind}, Received: val}
	}

	switch vr.Type() {
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	default:
		return fmt.Errorf("cannot decode %v into a %s", vr.Type(), val.Type())
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	if val.IsNil() {
		val.Set(reflect.MakeMap(val.Type()))
	}

	valueType := val.Type().Elem()
	decoder, err := dctx.LookupDecoder(valueType)
	if err != nil {
		return err
	}

	if valueType == reflect.TypeOf((*interface{})(nil)).Elem() {
		dctx.Ancestor = val.Type()
	}

	for {
		key, elemR, err := dr.ReadElement()
		if err == bsonrw.ErrEOD {
			break
		}
		if err != nil {
			return err
		}

		elem := reflect.New(valueType).Elem()

		err = decoder.DecodeValue(dctx, elemR, elem)
		if err != nil {
			return err
		}

		typedKey, err := c.keyDecoder(key)
		if err != nil {
			return err
		}

		val.SetMapIndex(typedKey, elem)
	}
	return nil
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
