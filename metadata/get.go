package metadata

import (
	"kproxy/helpers"
	"os"
	"strconv"
	"time"
)

func GetExpired(fileName string) bool {
	db := getDatabaseSingleton()
	value, err := db.Get([]byte(fileName + "-expiry"))
	if err != nil || value == nil {
		return false
	}

	numericValue, err := strconv.Atoi(string(value))
	if err != nil {
		return false
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

func GetVisits(fileName string) int {
	db := getDatabaseSingleton()
	value, err := db.Get([]byte(fileName + "-visits"))
	if err != nil || value == nil {
		return 0
	}

	numericValue, err := strconv.Atoi(string(value))
	if err != nil {
		return 0
	}

	return numericValue
}

func GetStat(fileName string) os.FileInfo {
	file, err := os.Stat(helpers.GetObjectPath(fileName))
	if err != nil {
		return nil
	}

	return file
}
