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

func SetRelevantHeaders(fileName string, header, clientHeader http.Header, headerNames []string) {
	resource := Get(fileName)
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

	fullServerChecksum := ServerUrlSum(fileName, clientHeader, header)
	if fullServerChecksum != stringToSum(fileName) {
		specificResource := Get(fullServerChecksum)
		specificResource.Headers = resource.Headers
		specificResource.RequestHeaders = clientHeader
		specificResource.Save()
	}
}
