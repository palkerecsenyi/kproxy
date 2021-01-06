package metadata

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _collection *mongo.Collection
var _context context.Context
var _cancel context.CancelFunc

func getCollectionSingleton() (*mongo.Collection, context.Context) {
	_context, _cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer _cancel()
	return _collection, _context
}

func Init() {
	_context, _cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer _cancel()

	client, err := mongo.Connect(_context, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic("Could not connect!")
	}

	defer func() {
		if err = client.Disconnect(_context); err != nil {
			panic(err)
		}
	}()

	_collection = client.Database("kproxy").Collection("file-metadata")
}
