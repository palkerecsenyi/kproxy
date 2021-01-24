package metadata

import (
	"kproxy/helpers"
	"net/http"
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

func SetRelevantHeaders(url string, header, clientHeader http.Header, headerNames []string) {
	resource := Get(stringToSum(url))
	resource.RequestHeaders = clientHeader
	resource.Headers = make(http.Header)

	for key, values := range header {
		if !helpers.SliceIterator(func(value string) bool {
			return http.CanonicalHeaderKey(key) == http.CanonicalHeaderKey(value)
		}, headerNames) {
			continue
		}

		for _, value := range values {
			resource.Headers.Add(key, value)
		}
	}

	resource.Save()

	fullServerChecksum := ServerUrlSum(url, clientHeader, header)
	if fullServerChecksum != stringToSum(url) {
		specificResource := Get(fullServerChecksum)
		specificResource.Headers = resource.Headers.Clone()
		specificResource.RequestHeaders = clientHeader
		specificResource.Save()
	}
}

func SetForceCache(fileName string, forced bool) {
	resource := Get(fileName)
	resource.CachedForOverride = forced
	resource.Save()
}
