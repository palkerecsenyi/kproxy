package cache

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"os"
)

func Get(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	if !shouldGetFromCache(req) {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	urlSum := getUrlSum(ctx)
	data, err := os.ReadFile(getObjectPath(urlSum))
	if err != nil {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	response := goproxy.NewResponse(
		req,
		goproxy.ContentTypeHtml,
		http.StatusOK,
		string(data),
	)

	response.Header.Add("X-Proxy-Sum", urlSum)
	response.Header.Set("Cache-Control", "no-cache")

	ctx.UserData = ProxyCacheState{
		FromCache: true,
	}

	return response
}
