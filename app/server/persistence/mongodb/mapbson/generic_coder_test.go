package mapbson

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

type customInt int

func customIntEncoder(key reflect.Value) (string, error) {
	gID := key.Interface().(customInt)
	return strconv.Itoa(int(gID)), nil
}
func customIntDecoder(key string) (reflect.Value, error) {
	i, err := strconv.Atoi(key)
	if err != nil {
		return reflect.Value{}, err
	}
	typedKey := customInt(i)
	return reflect.ValueOf(typedKey), nil
}

type hasCustomInt struct {
	CustomMap map[customInt]string `bson:"customMap"`
}

func getCustomRegistry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()

	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)

	ciMapType := reflect.TypeOf(hasCustomInt{}.CustomMap)

	cmc := newCustomMapCoder(
		ciMapType,
		customIntEncoder,
		customIntDecoder,
	)

	rb.RegisterEncoder(ciMapType, cmc)
	rb.RegisterDecoder(ciMapType, cmc)

	return rb.Build()
}

func TestCustomMapCoder(t *testing.T) {
	registry := getCustomRegistry()

	testCases := []struct {
		msg       string
		input     hasCustomInt
		expOutput hasCustomInt
	}{{
		msg: `custom map with three entries`,
		input: hasCustomInt{
			CustomMap: map[customInt]string{
				customInt(1): `one`,
				customInt(2): `two`,
				customInt(3): `three`,
			},
		},
		expOutput: hasCustomInt{
			CustomMap: map[customInt]string{
				customInt(1): `one`,
				customInt(2): `two`,
				customInt(3): `three`,
			},
		},
	}, {
		msg: `empty map`,
		input: hasCustomInt{
			CustomMap: map[customInt]string{},
		},
		expOutput: hasCustomInt{
			CustomMap: map[customInt]string{},
		},
	}}

	for _, tc := range testCases {
		data, err := bson.MarshalWithRegistry(registry, tc.input)
		require.NoError(t, err, tc.msg)

		actOutput := hasCustomInt{}
		err = bson.UnmarshalWithRegistry(registry, data, &actOutput)
		require.NoError(t, err, tc.msg)
		assert.Equal(t, tc.expOutput, actOutput, tc.msg)
	}
}
