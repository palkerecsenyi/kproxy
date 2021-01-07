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

	urlSum := helpers.GetUrlSum(ctx)
	maxAge := helpers.GetMaxAge(resp)
	metadata.SetMaxAge(urlSum, maxAge)

	body := helpers.ResponseToBytes(resp)
	if body == nil {
		return
	}

	contentType := helpers.GetMimeTypeFromHeader(resp)
	metadata.SetMimeType(urlSum, contentType)

	err := os.WriteFile(
		helpers.GetObjectPath(urlSum),
		body,
		0777,
	)

	if err != nil {
		fmt.Println(err)
	}
}
