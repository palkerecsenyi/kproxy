package config

import (
	"crypto/rand"
	"net/http"
)

func getSpeedTestPage(res http.ResponseWriter, _ *http.Request) {
	writeTemplate("speedtest", nil, res)
}

func startSpeedTest(res http.ResponseWriter, _ *http.Request) {
	randomData := make([]byte, 250000000)
	_, err := rand.Read(randomData)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	res.Header().Add("Cache-Control", "no-cache")
	res.Header().Add("Cache-Control", "no-store")
	_, _ = res.Write(randomData)
}
