package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func hasCollectionIndex(ctx context.Context, idxs mongo.IndexView, indexName string) (bool, error) {
	opt := &options.ListIndexesOptions{}
	opt.SetMaxTime(5 * time.Second)
	cur, err := idxs.List(ctx, opt)
	if err != nil {
		return false, err
	}

	for cur.Next(ctx) {
		index := bson.D{}
		err = cur.Decode(&index)
		if err != nil {
			return false, err
		}
		for _, i := range index {
			if key := i.Key; key == `key` {
				if val, ok := i.Value.(bson.D); ok && val[0].Key == indexName {
					// found the desired index
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func createCollectionIndex(ctx context.Context, idxs mongo.IndexView, indexName string) error {
	// Using `1` means "create ascending index"
	// See https://docs.mongodb.com/manual/reference/method/db.collection.createIndex/
	keys := bsonx.Doc{{
		Key:   indexName,
		Value: bsonx.Int64(int64(1)),
	}}
	im := mongo.IndexModel{}
	im.Keys = keys
	opts := options.CreateIndexes().SetMaxTime(5 * time.Second)

	_, err := idxs.CreateOne(ctx, im, opts)
	return err
}
