package config

import (
	"crypto/rand"
	"net/http"
	"strconv"
	"time"
)

func getSpeedTestPage(res http.ResponseWriter, _ *http.Request) {
	writeTemplate("speedtest", nil, res)
}

func startSpeedTest(res http.ResponseWriter, req *http.Request) {
	megabytesString := req.URL.Query().Get("mb")
	megabytes, err := strconv.ParseFloat(megabytesString, 64)
	if megabytes <= 0 || megabytes > 500 || err != nil {
		res.WriteHeader(400)
		return
	}

	start := time.Now()
	bytes := int(megabytes * 1000000)
	randomData := make([]byte, bytes)
	_, err = rand.Read(randomData)
	if err != nil {
		res.WriteHeader(500)
		return
	}
	elapsed := time.Since(start)

	res.Header().Add("Content-Length", strconv.Itoa(bytes))
	res.Header().Add("Cache-Control", "no-cache, no-store, private")
	res.Header().Add("Server-Timing", "gen;dur="+strconv.FormatInt(elapsed.Milliseconds(), 10))
	_, _ = res.Write(randomData)

	if f, ok := res.(http.Flusher); ok {
		f.Flush()
	}
}
