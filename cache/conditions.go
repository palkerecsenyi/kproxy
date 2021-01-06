package cache

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"strings"
)

type ProxyCacheState struct {
	// true if the requested resource was taken from cache
	// if true, no request was made
	FromCache bool
}

func _sliceIterator(iterator func(value string) bool, slice []string) bool {
	for _, value := range slice {
		if iterator(value) {
			return true
		}
	}
	return false
}

func _sliceContainsPrefix(searchValue string, slice []string) bool {
	return _sliceIterator(func(value string) bool {
		return strings.HasPrefix(searchValue, value)
	}, slice)
}

func _sliceContainsString(searchValue string, slice []string) bool {
	return _sliceIterator(func(value string) bool {
		return value == searchValue
	}, slice)
}

func _headerContainsAny(headerValue string, key ...string) bool {
	headerSlice := strings.Split(headerValue, ", ")
	return _sliceIterator(func(value string) bool {
		return _sliceContainsString(value, headerSlice)
	}, key)
}

func shouldSave(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
	// don't cache if the server doesn't want us to
	cacheControl := resp.Header.Get("Cache-Control")
	if _headerContainsAny(cacheControl, "no-cache", "no-store", "private") {
		return false
	}

	// don't even _try_ to cache anything that isn't http/s
	urlScheme := ctx.Req.URL.Scheme
	if urlScheme != "http" && urlScheme != "https" {
		return false
	}

	// other methods are almost never repeatable
	method := ctx.Req.Method
	if method != "GET" {
		return false
	}

	// no point in saving unsuccessful requests or redirects
	// however, 304s are often sent by other caching servers, so we can cache those too
	responseCode := resp.StatusCode
	if responseCode < 200 || (responseCode >= 300 && responseCode != 304) {
		return false
	}

	// we'll interpolate binary types later during saving/getting
	// this list is just a general filter
	contentType := resp.Header.Get("Content-Type")
	allowedContentTypes := []string{
		"text/html",
		"text/css",
		"application/javascript",
		"text/javascript", // for compatibility with bad servers and websites
		"image/png",
		"image/jpeg",
		"image/x-icon", // an unofficial type primarily used by .ico files, which almost all websites use (favicon.ico)
	}
	if !_sliceContainsPrefix(contentType, allowedContentTypes) {
		return false
	}

	return true
}

func shouldGetFromCache(req *http.Request) bool {
	// don't cache if the client doesn't want us to
	// e.g. when pressing ctrl/cmd+shift+r in chrome, it will user a Cache-Control: no-cache header
	cacheControl := req.Header.Get("Cache-Control")
	if strings.Contains(cacheControl, "no-cache") {
		return false
	}

	return true
}
