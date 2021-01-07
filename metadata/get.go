package metadata

import (
	"strconv"
	"time"
)

func GetExpired(fileName string) bool {
	db := getDatabaseSingleton()
	value, err := db.Get([]byte(fileName + "-expiry"))
	if err != nil || value == nil {
		return true
	}

	numericValue, err := strconv.Atoi(string(value))
	if err != nil {
		return true
	}

	expiry := time.Unix(int64(numericValue), 0)
	return expiry.Before(time.Now())
}

func GetMimeType(fileName string) string {
	db := getDatabaseSingleton()
	value, err := db.Get([]byte(fileName + "-mime"))
	if err != nil || value == nil {
		return ""
	}

	return string(value)
}
