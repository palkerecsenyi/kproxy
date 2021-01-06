package cache

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"kproxy/metadata"
	"net/http"
	"os"
)

func Save(resp *http.Response, ctx *goproxy.ProxyCtx) {
	if !shouldSave(resp, ctx) {
		return
	}

	urlSum := getUrlSum(ctx)
	maxAge := getMaxAge(resp)
	metadata.SetMaxAge(urlSum, maxAge)

	body := responseToBytes(resp)
	if body == nil {
		return
	}

	err := os.WriteFile(
		getObjectPath(urlSum),
		body,
		0777,
	)

	if err != nil {
		fmt.Println(err)
	}
}
