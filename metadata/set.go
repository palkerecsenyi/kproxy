package metadata

import (
	"kproxy/helpers"
	"net/http"
	"time"
)

type MultiOperationContext struct {
	PrimaryResource  string
	SpecificResource string
}

// generates a context to perform operations through
func MultiOperation(
	url string,
	responseHeader,
	requestHeader http.Header,
) MultiOperationContext {
	context := MultiOperationContext{
		PrimaryResource: stringToSum(url),
	}
	context.SpecificResource = ServerUrlSum(url, requestHeader, responseHeader)

	return context
}

func SingleOperation(url string) MultiOperationContext {
	summedUrl := stringToSum(url)
	return MultiOperationContext{
		PrimaryResource:  summedUrl,
		SpecificResource: summedUrl,
	}
}

func (context *MultiOperationContext) performOperation(callback func(resourceId string)) {
	callback(context.PrimaryResource)
	if context.PrimaryResource != context.SpecificResource {
		callback(context.SpecificResource)
	}
}

func (context *MultiOperationContext) SetMaxAge(maxAge time.Duration) {
	if maxAge.Seconds() < 0 {
		return
	}

	context.performOperation(func(resourceId string) {
		resource := Get(resourceId)
		resource.Expiry = time.Now().Add(maxAge).Unix()
		resource.Save()
	})
}

func (context *MultiOperationContext) SetMimeType(mimeType string) {
	context.performOperation(func(resourceId string) {
		resource := Get(resourceId)
		resource.MimeType = mimeType
		resource.Save()
	})
}

func (context *MultiOperationContext) IncrementVisits() {
	context.performOperation(func(resourceId string) {
		resource := Get(resourceId)
		resource.IncrementVisits()
	})
}

func (context *MultiOperationContext) SetRelevantHeaders(headerNames []string, responseHeader, requestHeader http.Header) {
	context.performOperation(func(resourceId string) {
		resource := Get(resourceId)
		resource.RequestHeaders = requestHeader
		resource.Headers = make(http.Header)

		for key, values := range responseHeader {
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
	})
}

// set whether the resource is only being cached because of an override
func (context *MultiOperationContext) SetForceCache(forced bool) {
	context.performOperation(func(resourceId string) {
		resource := Get(resourceId)
		resource.CachedForOverride = forced
		resource.Save()
	})
}
