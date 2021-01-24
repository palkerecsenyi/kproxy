package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"kproxy/cache"
	"kproxy/certificate"
	"kproxy/config"
	"kproxy/cron"
	"kproxy/metadata"
	"log"
	"net/http"
	"regexp"
)

func main() {
	cleanMode := flag.Bool("clean", false, "Run a cache clean.")
	port := flag.String("port", "80", "The port to run the proxy server on")
	enableConfigServer := flag.Bool("config", true, "Enable an HTTP server to allow proxy configuration")
	configServerPort := flag.String("config-port", "8080", "The port to run the HTTP config server on (if enabled)")
	flag.Parse()

	metadata.Init()
	if *cleanMode {
		fmt.Println("Running a cache clean!")
		cron.Clean()
		return
	}

	if *enableConfigServer {
		if *configServerPort == "" {
			panic("No port given for config server.")
		}
		config.Start(*configServerPort)
	}

	certificate.SetCA()

	proxyServer := goproxy.NewProxyHttpServer()

	condition := goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))
	proxyServer.OnRequest(condition).HandleConnect(goproxy.AlwaysMitm)

	proxyServer.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, cache.Get(req, ctx)
	})

	proxyServer.OnResponse(condition).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp.StatusCode >= 400 {
			response := goproxy.TextResponse(ctx.Req, "kProxy: "+resp.Status)
			response.StatusCode = resp.StatusCode
			response.Status = resp.Status
			return response
		}

		if userData, ok := ctx.UserData.(cache.ProxyCacheState); ok {

			resp.Header.Add("X-Cache-User", metadata.GetUserId(ctx.Req))
			if userData.FromCache {
				resp.Header.Add("X-Cache", "Hit from kProxy")
			} else {
				resp.Header.Add("X-Cache-Sum", metadata.ServerUrlSum(
					resp.Request.URL.String(),
					resp.Request.Header,
					resp.Header,
				))

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
