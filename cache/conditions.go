package cache

import (
	"github.com/elazarl/goproxy"
	"kproxy/helpers"
	"net/http"
	"strings"
)

type ProxyCacheState struct {
	// true if the requested resource was taken from cache
	// if true, no request was made
	FromCache bool
}

func _headerContainsAny(headerValue string, key ...string) bool {
	headerSlice := strings.Split(headerValue, ", ")
	return helpers.SliceContainsAnyString(headerSlice, key...)
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

	// we'll interpolate binary types later using helpers.IsTextualMime
	// this list is just a general filter
	contentType := resp.Header.Get("Content-Type")
	if !helpers.SliceContainsPrefix(contentType, allowedContentTypes) {
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

	// get and save may be different is a URL supports multiple methods
	method := req.Method
	if method != "GET" {
		return false
	}

	return true
}
