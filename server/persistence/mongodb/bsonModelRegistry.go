package mongodb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/joshprzybyszewski/cribbage/model"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func modelBSONRegistry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()

	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	registerGameIDCoders(rb)
	registerPlayerColorCoders(rb)

	return rb.Build()
}

func registerGameIDCoders(rb *bsoncodec.RegistryBuilder) {
	var refTypeGameID reflect.Type = reflect.TypeOf(model.InvalidGameID)
	c := &gameIDCoder{
		gameIDKind: refTypeGameID.Kind(),
	}

	// rb.RegisterEncoder(refTypeGameID, c)
	rb.RegisterDecoder(refTypeGameID, c)
}

func registerPlayerColorCoders(rb *bsoncodec.RegistryBuilder) {
	var refPlayerColorType reflect.Type = reflect.TypeOf(model.Green)
	c := &playerColorCoder{
		playerColorKind: refPlayerColorType.Kind(),
	}

	// rb.RegisterEncoder(refPlayerColorType, c)
	rb.RegisterDecoder(refPlayerColorType, c)
}

var _ bsoncodec.ValueEncoder = (*gameIDCoder)(nil)
var _ bsoncodec.ValueDecoder = (*gameIDCoder)(nil)

type gameIDCoder struct {
	gameIDKind reflect.Kind // should be reflect.TypeOf(model.InvalidGameID).Kind()
}

func (c *gameIDCoder) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Kind() != c.gameIDKind {
		return bsoncodec.ValueEncoderError{Name: "GameIDEncodeValue", Kinds: []reflect.Kind{c.gameIDKind}, Received: val}
	}

	vi := val.String()
	fmt.Printf("game id: write val.String() := %+v\n", val.String())
	return vw.WriteString(vi)
}

func (c *gameIDCoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() != bsontype.String {
		return fmt.Errorf("cannot decode %v into a GameID", vr.Type())
	}
	if !val.IsValid() || !val.CanSet() || val.Kind() != c.gameIDKind {
		return bsoncodec.ValueDecoderError{Name: "GameIDDecodeValue", Kinds: []reflect.Kind{c.gameIDKind}, Received: val}
	}

	s, err := vr.ReadString()
	if err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	gID := model.GameID(i)

	val.Set(reflect.ValueOf(gID))
	return nil
}

var _ bsoncodec.ValueEncoder = (*playerColorCoder)(nil)
var _ bsoncodec.ValueDecoder = (*playerColorCoder)(nil)

type playerColorCoder struct {
	playerColorKind reflect.Kind // should be reflect.TypeOf(model.InvalidGameID).Kind()
}

func (c *playerColorCoder) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Kind() != c.playerColorKind {
		return bsoncodec.ValueEncoderError{Name: "GameIDEncodeValue", Kinds: []reflect.Kind{c.playerColorKind}, Received: val}
	}

	s := val.String()
	fmt.Printf("pc encode: %+v\n", s)
	return vw.WriteString(s)
}

func (c *playerColorCoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != c.playerColorKind {
		return bsoncodec.ValueDecoderError{
			Name:     "PlayerColor.DecodeValue",
			Kinds:    []reflect.Kind{c.playerColorKind},
			Received: val,
		}
	}

	if vr.Type() == bsontype.String {
		s, err := vr.ReadString()
		if err != nil {
			return err
		}

		fmt.Printf("string %+v\n", s)
		for i := model.Green; i < 5; i++ {
			if i.String() == s {
				val.Set(reflect.ValueOf(i))
				return nil
			}
		}
		return errors.New(`did not find playercolor`)
	} else if vr.Type() == bsontype.Int32 {
		i, err := vr.ReadInt32()
		if err != nil {
			return err
		}

		fmt.Printf("int32 %+v\n", i)
		pc := model.PlayerColor(i)
		val.Set(reflect.ValueOf(pc))
		return nil
	}

	return fmt.Errorf("cannot decode %v into a PlayerColor", vr.Type())
}
