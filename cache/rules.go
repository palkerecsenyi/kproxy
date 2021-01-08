package cache

import (
	"github.com/gobwas/glob"
	"strings"
)

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
	"application/font-woff2",
	"font/otf",

	"audio/mpeg",
	"video/mp4",
	"video/mpeg",
}

type cacheRule struct {
	glob      glob.Glob
	onlyTypes []string
}

var alwaysCache = []cacheRule{
	{
		glob: glob.MustCompile("*.wikipedia.org"),
		onlyTypes: []string{
			"text/html",
		},
	},
}

var neverCache = []cacheRule{
	{
		glob: glob.MustCompile("example.net"),
	},
}

// returns true for a positive match, false for no match
func testRule(ruleSlice []cacheRule, url, contentType string) bool {
	for _, item := range ruleSlice {
		urlMatchesGlob := item.glob.Match(url)
		if !urlMatchesGlob {
			continue
		}

		if item.onlyTypes != nil {
			onlyTypeMatched := false
			for _, onlyType := range item.onlyTypes {
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
func shouldCacheUrl(url, contentType string) int {
	// neverCache rules take priority over alwaysCache
	if testRule(neverCache, url, contentType) {
		return forceNoCache
	}

	if testRule(alwaysCache, url, contentType) {
		return forceCache
	}

	return noRule
}
