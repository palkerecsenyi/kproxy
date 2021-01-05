package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hello, world! this is definitely a proxy server")
	proxyServer := goproxy.NewProxyHttpServer()
	err := http.ListenAndServe(":8080", proxyServer)

	if err != nil {
		log.Fatal(err)
	}
}
