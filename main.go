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
	"kproxy/metadata/analytics"
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
	proxyServer.Tr = &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		ForceAttemptHTTP2: true,
	}

	condition := goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))
	proxyServer.OnRequest(condition).HandleConnect(goproxy.AlwaysMitm)

	proxyServer.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, cache.Get(req, ctx)
	})

	proxyServer.OnResponse(condition).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if userData, ok := ctx.UserData.(cache.ProxyCacheState); ok {

			resourceOperations := metadata.MultiOperation(resp.Request.URL.String(), resp.Header, userData.RequestHeaders)
			resourceOperations.IncrementVisits()

			resp.Header.Add("X-Cache-User", metadata.GetUserId(ctx.Req))
			if userData.FromCache {
				go analytics.LogRequest(resp.Request.URL, true, uint64(resp.ContentLength))
				resp.Header.Add("X-Cache", "Hit from kProxy")
			} else {
				go analytics.LogRequest(resp.Request.URL, false, 0)
				resp.Header.Add("X-Cache-Sum", resourceOperations.SpecificResource)
				resp.Header.Add("X-Cache", "Miss from kProxy")
				defer cache.Save(resp, ctx)
			}
		}

		return resp
	})

	log.Println("Listening on " + *port)
	err := http.ListenAndServe(":"+*port, proxyServer)
	if err != nil {
		log.Fatal(err)
	}
}
