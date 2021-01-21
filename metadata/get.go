package metadata

import (
	"kproxy/helpers"
	"os"
	"strconv"
	"time"
)

// returns: (expired), (expires in seconds — 0 if expired)
func GetExpired(fileName string) (bool, int) {
	db := GetDatabaseSingleton()
	value, err := db.Get([]byte(fileName + "-expiry"))
	if err != nil || value == nil {
		return true, 0
	}

	numericValue, err := strconv.Atoi(string(value))
	if err != nil {
		return true, 0
	}

	expiry := time.Unix(int64(numericValue), 0)
	expired := expiry.Before(time.Now())
	if expired {
		return true, 0
	} else {
		return false, int(expiry.Sub(time.Now()).Seconds())
	}
}

func GetMimeType(fileName string) string {
	db := GetDatabaseSingleton()
	value, err := db.Get([]byte(fileName + "-mime"))
	if err != nil || value == nil {
		return ""
	}

	return string(value)
}

func GetVisits(fileName string) int {
	db := GetDatabaseSingleton()
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
