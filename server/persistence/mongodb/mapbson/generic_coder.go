package mapbson

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func newCustomMapCoder(
	mapType reflect.Type,
	keyEnc func(reflect.Value) (string, error),
	keyDec func(string) (reflect.Value, error),
) *customMapCoder {
	return &customMapCoder{
		mapKind:    mapType.Kind(),
		keyKind:    mapType.Key().Kind(),
		valueKind:  mapType.Elem().Kind(),
		keyEncoder: keyEnc,
		keyDecoder: keyDec,
	}
}

var _ bsoncodec.ValueEncoder = (*customMapCoder)(nil)
var _ bsoncodec.ValueDecoder = (*customMapCoder)(nil)

type customMapCoder struct {
	mapKind   reflect.Kind
	keyKind   reflect.Kind
	valueKind reflect.Kind

	keyEncoder func(reflect.Value) (string, error)
	keyDecoder func(string) (reflect.Value, error)
}

func (c *customMapCoder) EncodeValue(
	ectx bsoncodec.EncodeContext,
	vw bsonrw.ValueWriter,
	val reflect.Value,
) error {

	if !val.IsValid() ||
		val.Kind() != reflect.Map ||
		val.Type().Key().Kind() != c.keyKind ||
		val.Type().Elem().Kind() != c.valueKind {
		return bsoncodec.ValueEncoderError{
			Name:     "CustomMapCoderEncodeValue",
			Kinds:    []reflect.Kind{c.mapKind},
			Received: val,
		}
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

func (c *customMapCoder) DecodeValue(
	dctx bsoncodec.DecodeContext,
	vr bsonrw.ValueReader,
	val reflect.Value,
) error {

	if !val.CanSet() ||
		val.Kind() != reflect.Map ||
		val.Type().Key().Kind() != c.keyKind ||
		val.Type().Elem().Kind() != c.valueKind {
		return bsoncodec.ValueDecoderError{
			Name:     "CustomMapCoderDecodeValue",
			Kinds:    []reflect.Kind{c.mapKind},
			Received: val,
		}
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
