package metadata

import (
	"kproxy/helpers"
	"net/http"
	"strings"
	"time"
)

func SetMaxAge(fileName string, maxAge time.Duration) {
	if maxAge.Seconds() < 0 {
		return
	}

	resource := Get(fileName)
	resource.Expiry = time.Now().Add(maxAge).Unix()
	resource.Save()
}

func SetMimeType(fileName, mimeType string) {
	resource := Get(fileName)
	resource.MimeType = mimeType
	resource.Save()
}

func IncrementVisits(fileName string) {
	resource := Get(fileName)
	resource.IncrementVisits()
}

func SetRelevantHeaders(fileName string, header http.Header, headerNames []string) {
	resource := Get(fileName)
	resource.Headers = make(http.Header)

	for key, values := range header {
		if !helpers.SliceContainsString(strings.ToLower(key), headerNames) {
			continue
		}

		for _, value := range values {
			resource.Headers.Add(key, value)
		}
	}

	resource.Save()
}
