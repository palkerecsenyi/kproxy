package config

import (
	"log"
	"net/http"
)

func Start(port string) {
	server := http.NewServeMux()

	server.HandleFunc("/settings", getSettings)
	server.HandleFunc("/settings/save", saveSettings)
	server.HandleFunc("/settings/delete", deleteCacheRule)

	server.HandleFunc("/schedule-download", downloadLargeFile)
	server.HandleFunc("/download-status", downloadStatus)
	server.HandleFunc("/download-file", downloadSavedFile)
	server.HandleFunc("/speed-test", getSpeedTestPage)
	server.HandleFunc("/speed-test/random", startSpeedTest)

	server.HandleFunc("/kproxy.pem", downloadCert)
	server.HandleFunc("/test", testCache)
	server.HandleFunc("/logs", getLogs)
	server.HandleFunc("/", reportStatus)

	go func() {
		log.Println("Starting config server on " + port)
		_ = http.ListenAndServe(":"+port, server)
	}()
}
