package cache

import (
	"github.com/elazarl/goproxy"
	"kproxy/helpers"
	"kproxy/metadata"
	"net/http"
	"strings"
)

type ProxyCacheState struct {
	// true if the requested resource was taken from cache
	// if true, no request was made
	FromCache      bool
	RequestHeaders http.Header
}

func _headerContainsAny(headerValue string, key ...string) bool {
	headerSlice := strings.Split(headerValue, ", ")
	return helpers.SliceContainsAnyString(headerSlice, key...)
}

// RFC7234 Section 3: Storing responses in Caches https://tools.ietf.org/html/rfc7234#section-3
func shouldSave(resp *http.Response, ctx *goproxy.ProxyCtx, urlSum string) bool {
	// point 1 – requests that aren't http or https are not understood by the server
	urlScheme := ctx.Req.URL.Scheme
	if urlScheme != "http" && urlScheme != "https" {
		return false
	}

	// point 2 – only successful responses
	responseCode := resp.StatusCode
	if responseCode < 200 || responseCode > 299 {
		return false
	}

	// also point 1 — non-GET requests are buggy to cache
	method := ctx.Req.Method
	if method != "GET" {
		return false
	}

	// point 1 — we only cache MIME types that are configured as cacheable
	// responses without Content-Type headers are inherently invalid
	contentType := resp.Header.Get("Content-Type")
	if !helpers.SliceContainsPrefix(contentType, allowedContentTypes) {
		return false
	}

	// point 5 — cannot be overridden for security reasons
	if resp.Header.Get("Authorization") != "" {
		// does not implement section 3.2 yet
		return false
	}

	// these overrides only override server Cache-Control headers
	// configurable overrides aren't defined in RFC7234
	switch shouldCacheUrl(ctx.Req, contentType) {
	case forceCache:
		return true
	case forceNoCache:
		return false
	case noRule:
		if metadata.GetForceCache(urlSum) {
			return false
		}
	}

	// point 2 & 3
	// currently, we treat no-cache as synonymous to no-store
	cacheControl := resp.Header.Get("Cache-Control")
	if _headerContainsAny(cacheControl, "no-cache", "no-store", "private") {
		return false
	}

	return true
}

func shouldGetFromCache(req *http.Request, contentType, urlSum string) bool {
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

	if override := shouldCacheUrl(req, contentType); override == forceNoCache {
		return false
	} else if override == noRule && metadata.GetForceCache(urlSum) {
		return false
	}

	return true
}
