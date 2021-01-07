package metadata

import "go.mongodb.org/mongo-driver/bson"

const DefaultType = "text/html"

func GetMimeType(fileName string) string {
	collection, ctx := getCollectionSingleton()
	if collection == nil {
		return DefaultType
	}

	result := collection.FindOne(ctx, bson.M{
		"name": fileName,
	})

	if result.Err() != nil {
		return DefaultType
	}

	data := DocumentData{}
	err := result.Decode(data)
	if err != nil {
		return DefaultType
	}

	if data.mimeType == "" {
		return DefaultType
	}

	return data.mimeType
}
