package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"kproxy/cache"
	"kproxy/certificate"
	"kproxy/metadata"
	"log"
	"net/http"
	"regexp"
)

func main() {
	cleanMode := flag.Bool("clean", false, "Run a cache clean.")
	port := flag.String("port", "80", "The port to run the proxy server on")
	flag.Parse()

	metadata.Init()
	if *cleanMode {
		fmt.Println("Running a cache clean!")
		metadata.Clean()
		return
	}

	certificate.SetCA()

	proxyServer := goproxy.NewProxyHttpServer()

	condition := goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))
	proxyServer.OnRequest(condition).HandleConnect(goproxy.AlwaysMitm)

	proxyServer.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, cache.Get(req, ctx)
	})

	proxyServer.OnResponse(condition).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if userData, ok := ctx.UserData.(cache.ProxyCacheState); ok {
			if userData.FromCache {
				resp.Header.Add("X-Cache", "Hit from kProxy")
			} else {
				resp.Header.Add("X-Cache", "Miss from kProxy")
				defer cache.Save(resp, ctx)
			}
		}

		return resp
	})

	err := http.ListenAndServe(":"+*port, proxyServer)

	if err != nil {
		log.Fatal(err)
	}
}
