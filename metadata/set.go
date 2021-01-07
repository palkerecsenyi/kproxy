package metadata

import (
	"strconv"
	"time"
)

func SetMaxAge(fileName string, maxAge time.Duration) {
	if maxAge.Nanoseconds() == 0 {
		return
	}

	db := getDatabaseSingleton()
	expiry := strconv.FormatInt(time.Now().Add(maxAge).Unix(), 10)
	_ = db.Put([]byte(fileName+"-expiry"), []byte(expiry))
}

func SetMimeType(fileName, mimeType string) {
	db := getDatabaseSingleton()
	_ = db.Put([]byte(fileName+"-mime"), []byte(mimeType))
}

func IncrementVisits(fileName string) {
	visits := GetVisits(fileName)
	visits++

	db := getDatabaseSingleton()
	_ = db.Put([]byte(fileName+"-visits"), []byte(strconv.Itoa(visits)))
}
