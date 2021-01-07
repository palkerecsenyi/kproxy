package metadata

import (
	"go.mongodb.org/mongo-driver/bson"
	"kproxy/helpers"
	"os"
	"path"
	"time"
)

func Clean() {
	collection, ctx := getCollectionSingleton()
	contents, err := os.ReadDir(helpers.GetPath())
	if err != nil {
		panic(err)
	}

	for _, file := range contents {
		fileName := file.Name()
		filePath := path.Join(helpers.GetPath(), fileName)

		result := collection.FindOne(ctx, bson.M{
			"name": fileName,
		})
		if result.Err() != nil {
			continue
		}

		data := DocumentData{}
		err = result.Decode(data)
		if err != nil {
			continue
		}

		expiryDate := time.Unix(data.expiryDate, 0)
		if expiryDate.Before(time.Now()) {
			_ = os.Remove(filePath)
		}
	}
}
