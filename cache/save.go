package cache

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"kproxy/helpers"
	"kproxy/metadata"
	"net/http"
	"os"
)

func Save(resp *http.Response, ctx *goproxy.ProxyCtx) {
	urlSum := metadata.ServerUrlSum(resp.Request.URL.String(), resp.Request.Header, resp.Header)
	if !shouldSave(resp, ctx, urlSum) {
		return
	}

	contentType := helpers.GetMimeTypeFromHeader(resp)
	if shouldCacheUrl(ctx.Req, contentType) == forceCache {
		metadata.SetForceCache(urlSum, true)
	}

	maxAge := helpers.GetRequestMaxAge(resp)
	metadata.SetMaxAge(urlSum, maxAge)

	body := helpers.ResponseToBytes(resp)
	if body == nil {
		return
	}

	metadata.SetMimeType(urlSum, contentType)
	metadata.SetRelevantHeaders(resp.Request.URL.String(), resp.Header, resp.Request.Header, cacheableHeaders)

	err := os.WriteFile(
		helpers.GetObjectPath(urlSum),
		body,
		0777,
	)

	if err != nil {
		fmt.Println(err)
	}
}
