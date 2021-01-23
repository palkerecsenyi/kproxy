package metadata

import (
	"kproxy/helpers"
	"net/http"
	"os"
	"time"
)

// returns: (expired), (expires in seconds — 0 if expired)
func GetExpired(fileName string) (bool, int) {
	resource := Get(fileName)
	if resource.Expiry == 0 {
		return true, 0
	}

	expiry := time.Unix(resource.Expiry, 0)
	expired := expiry.Before(time.Now())
	if expired {
		return true, 0
	} else {
		return false, int(expiry.Sub(time.Now()).Seconds())
	}
}

func GetMimeType(fileName string) string {
	resource := Get(fileName)
	return resource.MimeType
}

func GetVisits(fileName string) int {
	resource := Get(fileName)
	return resource.Visits
}

func GetStat(fileName string) os.FileInfo {
	file, err := os.Stat(helpers.GetObjectPath(fileName))
	if err != nil {
		return nil
	}

	return file
}

func GetHeaders(fileName string) http.Header {
	header := make(http.Header)
	resource := Get(fileName)
	if resource.Headers == nil {
		return header
	}

	header = resource.Headers
	return header
}

func GetForceCache(fileName string) bool {
	resource := Get(fileName)
	return resource.CachedForOverride
}
