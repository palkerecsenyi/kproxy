package metadata

import (
	"go.mongodb.org/mongo-driver/bson"
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

	_, _ = collection.UpdateOne(ctx, bson.M{
		"name": fileName,
	}, bson.M{
		"$set": bson.M{
			"expiryDate": time.Now().Add(maxAge).Unix(),
		},
	}, upsertUpdate)
}

func SetMimeType(fileName, mimeType string) {
	collection, ctx := getCollectionSingleton()
	if collection == nil {
		return
	}

	_, _ = collection.UpdateOne(ctx, bson.M{
		"name": fileName,
	}, bson.M{
		"$set": bson.M{
			"mimeType": mimeType,
		},
	}, upsertUpdate)
}
