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
	userData, ok := ctx.UserData.(ProxyCacheState)
	if !ok {
		return
	}

	urlSum := metadata.ServerUrlSum(resp.Request.URL.String(), userData.RequestHeaders, resp.Header)
	if !shouldSave(resp, ctx, urlSum) {
		return
	}

	resourceOperations := metadata.MultiOperation(resp.Request.URL.String(), resp.Header, userData.RequestHeaders)
	contentType := helpers.GetMimeTypeFromHeader(resp)
	if shouldCacheUrl(resp.Request, contentType) == forceCache {
		resourceOperations.SetForceCache(true)
	}

	maxAge := helpers.GetRequestMaxAge(resp)
	resourceOperations.SetMaxAge(maxAge)

	body := helpers.ResponseToBytes(resp)
	if body == nil {
		return
	}

	resourceOperations.SetMimeType(contentType)
	resourceOperations.SetRelevantHeaders(cacheableHeaders, resp.Header, resp.Request.Header)

	err := os.WriteFile(
		helpers.GetObjectPath(urlSum),
		body,
		0777,
	)

	if err != nil {
		fmt.Println(err)
	}
}
