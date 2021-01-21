package config

import "net/http"

func Start(port string) {
	server := http.NewServeMux()

	server.HandleFunc("/schedule-download", downloadLargeFile)
	server.HandleFunc("/download-status", downloadStatus)
	server.HandleFunc("/download-file", downloadSavedFile)
	server.HandleFunc("/test", testCache)

	server.HandleFunc("/kproxy.pem", downloadCert)
	server.HandleFunc("/", reportStatus)

	go func() {
		_ = http.ListenAndServe(":"+port, server)
	}()
}
