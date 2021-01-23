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
	if !shouldSave(resp, ctx) {
		return
	}

	userData, ok := ctx.UserData.(ProxyCacheState)
	if !ok {
		return
	}

	urlSum := metadata.ServerUrlSum(ctx.Req.URL.String(), userData.RequestHeaders, resp.Header)
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
	metadata.SetRelevantHeaders(ctx.Req.URL.String(), resp.Header, userData.RequestHeaders, cacheableHeaders)

	err := os.WriteFile(
		helpers.GetObjectPath(urlSum),
		body,
		0777,
	)

	if err != nil {
		fmt.Println(err)
	}
}
