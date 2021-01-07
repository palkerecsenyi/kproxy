package metadata

import (
	"kproxy/helpers"
	"os"
	"path"
)

func Clean() {
	contents, err := os.ReadDir(helpers.GetPath())
	if err != nil {
		panic(err)
	}

	for _, file := range contents {
		fileName := file.Name()
		filePath := path.Join(helpers.GetPath(), fileName)

		if GetExpired(fileName) {
			_ = os.Remove(filePath)
		}
	}
}
