package metadata

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SetMaxAge(fileName string, maxAge time.Duration) {
	collection, ctx := getCollectionSingleton()
	if collection == nil {
		return
	}

	upsert := true
	_, _ = collection.UpdateOne(ctx, bson.M{
		"name": fileName,
	}, bson.M{
		"$set": bson.M{
			"expiryData": time.Now().Add(maxAge).Unix(),
		},
	}, &options.UpdateOptions{Upsert: &upsert})
}
