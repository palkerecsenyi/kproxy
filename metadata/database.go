package metadata

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _collection *mongo.Collection
var _context context.Context

func getCollectionSingleton() (*mongo.Collection, context.Context) {
	return _collection, _context
}

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_context = ctx
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	_collection = client.Database("kproxy").Collection("file-metadata")
}
