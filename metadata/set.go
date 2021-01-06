package metadata

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SetMaxAge(fileName string, maxAge time.Duration) {
	if maxAge.Nanoseconds() == 0 {
		return
	}

	collection, ctx := getCollectionSingleton()
	if collection == nil {
		return
	}

	upsert := true
	_, _ = collection.UpdateOne(ctx, bson.M{
		"name": fileName,
	}, bson.M{
		"$set": bson.M{
			"expiryDate": time.Now().Add(maxAge).Unix(),
		},
	}, &options.UpdateOptions{Upsert: &upsert})
}
