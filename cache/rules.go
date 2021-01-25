package cache

import (
	"github.com/gobwas/glob"
	"kproxy/metadata"
	"net/http"
	"strings"
)

// case-insensitive; uses http.CanonicalHeaderKey to standardise
// refers to response headers
var cacheableHeaders = []string{
	// access-control-allow-origin is set dynamically
	"access-control-allow-methods",
	"access-control-allow-credentials",
	"age",
	"expires",
	"vary",
	"last-modified",
}

// these MIME types are cached, and nothing else is
var allowedContentTypes = []string{
	"text/html",
	"text/css",
	"application/javascript",
	"text/javascript", // for compatibility with bad servers and websites

	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/webp",
	"image/gif",
	"image/svg+xml",
	"image/x-icon", // an unofficial type primarily used by .ico files, which almost all websites use (favicon.ico)

	"application/pdf",
	"font/ttf",
	"font/woff",
	"font/woff2",
	"application/font-woff",
	"application/font-woff2",
	"font/otf",

	"audio/mpeg",
	"video/mp4",
	"video/mpeg",
}

var alwaysCache = []metadata.CacheRule{
	{
		Glob: "*.wikipedia.org/*",
		OnlyTypes: []string{
			"text/html",
		},
	},
}

var neverCache = []metadata.CacheRule{
	{
		Glob: "**cloud.google.com/*",
	},
}

// returns true for a positive match, false for no match
func testRule(ruleSlice []metadata.CacheRule, url, contentType string) bool {
	for _, item := range ruleSlice {
		compiledGlob := glob.MustCompile(item.Glob)
		urlMatchesGlob := compiledGlob.Match(url)
		if !urlMatchesGlob {
			continue
		}

		if item.OnlyTypes != nil {
			onlyTypeMatched := false
			for _, onlyType := range item.OnlyTypes {
				if strings.HasPrefix(contentType, onlyType) {
					onlyTypeMatched = true
					break
				}
			}

			if onlyTypeMatched {
				return true
			} else {
				continue
			}
		} else {
			return true
		}
	}

	return false
}

const (
	noRule       = iota
	forceCache   = iota
	forceNoCache = iota
)

// returns one of the above constants
func shouldCacheUrl(req *http.Request, contentType string) int {
	simpleUrl := req.URL.Hostname() + req.URL.Path

	settings := metadata.GetSettings(req)
	if testRule(settings.NeverCache, simpleUrl, contentType) {
		return forceNoCache
	}
	if testRule(settings.AlwaysCache, simpleUrl, contentType) {
		return forceCache
	}

	// neverCache rules take priority over alwaysCache
	if testRule(neverCache, simpleUrl, contentType) {
		return forceNoCache
	}
	if testRule(alwaysCache, simpleUrl, contentType) {
		return forceCache
	}

	return noRule
}
