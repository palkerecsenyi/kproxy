package metadata

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var upsert = true
var upsertUpdate = &options.UpdateOptions{Upsert: &upsert}

type DocumentData struct {
	name       string
	expiryDate int64
	mimeType   string
}

var _collection *mongo.Collection
var _context context.Context

func getCollectionSingleton() (*mongo.Collection, context.Context) {
	_context, _ = context.WithTimeout(context.Background(), 5*time.Second)
	return _collection, _context
}

func Init() {
	_context, _ = context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(_context, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic("Could not connect!")
	}

	_collection = client.Database("kproxy").Collection("metadata")
}
