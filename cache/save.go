package cache

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"net/http"
	"os"
)

func Save(resp *http.Response, ctx *goproxy.ProxyCtx) {
	if !shouldSave(resp, ctx) {
		return
	}

	urlSum := getUrlSum(ctx)
	stringBody := responseToString(resp)

	err := os.WriteFile(
		getObjectPath(urlSum),
		[]byte(stringBody),
		0777,
	)

	if err != nil {
		fmt.Println(err)
	}
}
