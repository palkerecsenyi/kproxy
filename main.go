package main

import (
	"github.com/elazarl/goproxy"
	"kproxy/cache"
	"kproxy/certificate"
	"log"
	"net/http"
	"regexp"
)

func main() {
	certificate.SetCA()

	proxyServer := goproxy.NewProxyHttpServer()

	condition := goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))
	proxyServer.OnRequest(condition).HandleConnect(goproxy.AlwaysMitm)

	proxyServer.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, cache.Get(req, ctx)
	})

	proxyServer.OnResponse(condition).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp.Header.Add("X-Proxied-By", "kProxy")

		if userData, ok := ctx.UserData.(cache.ProxyCacheState); ok {
			if userData.FromCache {
				resp.Header.Add("X-Proxy-Source", "cache")
			} else {
				resp.Header.Add("X-Proxy-Source", "server")
			}
		}

		defer cache.Save(resp, ctx)
		return resp
	})

	err := http.ListenAndServe(":8080", proxyServer)

	if err != nil {
		log.Fatal(err)
	}
}
