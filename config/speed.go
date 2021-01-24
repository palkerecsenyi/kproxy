package config

import (
	"crypto/rand"
	"net/http"
	"strconv"
)

func getSpeedTestPage(res http.ResponseWriter, _ *http.Request) {
	writeTemplate("speedtest", nil, res)
}

func startSpeedTest(res http.ResponseWriter, req *http.Request) {
	megabytesString := req.URL.Query().Get("mb")
	megabytes, err := strconv.Atoi(megabytesString)
	if megabytes <= 0 || megabytes > 500 || err != nil {
		res.WriteHeader(400)
		return
	}

	randomData := make([]byte, megabytes*1000000)
	_, err = rand.Read(randomData)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	res.Header().Add("Cache-Control", "no-cache, no-store, private")
	_, _ = res.Write(randomData)

	if f, ok := res.(http.Flusher); ok {
		f.Flush()
	}
}
