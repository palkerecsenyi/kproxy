package config

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/dustin/go-humanize"
	"kproxy/certificate"
	"kproxy/eviction"
	"kproxy/helpers"
	"kproxy/metadata"
	"kproxy/metadata/analytics"
	"net/http"
	"os"
	"strconv"
	"time"
)

func reportStatus(res http.ResponseWriter, req *http.Request) {
	data := getJsonMap()

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		data["error"] = "Failed fetching CPU usage"
		sendJson(data, 500, res)
		return
	}

	cpuUsage := getJsonMap()
	cpuUsage["system"] = stat.CPUStatAll.System
	cpuUsage["user"] = stat.CPUStatAll.User
	cpuUsage["nice"] = stat.CPUStatAll.Nice
	data["cpu"] = cpuUsage

	storageUsage := eviction.CalculateStorageUsage()
	data["cache_usage_bytes"] = storageUsage
	data["cache_usage_human"] = humanize.Bytes(uint64(storageUsage))

	files, _ := os.ReadDir(helpers.GetPath())
	data["cache_object_count"] = len(files)

	data["your_ip"] = req.RemoteAddr
	data["my_time"] = time.Now().Format(time.RFC3339)

	sendJson(data, 200, res)
}

func getLogs(res http.ResponseWriter, req *http.Request) {
	data := getJsonMap()

	daysString := req.URL.Query().Get("days")
	days := 1
	if daysString != "" {
		days, _ = strconv.Atoi(daysString)
	}

	logs := analytics.GetLogs(time.Now().AddDate(0, 0, -days))
	totalSavings := analytics.SumSavings(logs)

	data["total_savings_bytes"] = totalSavings
	data["total_savings_human"] = humanize.Bytes(totalSavings)
	data["logs"] = logs

	sendJson(data, 200, res)
}

func downloadCert(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "application/x-pem-file")
	publicKey := certificate.GetPublicKey()
	_, _ = res.Write(publicKey)
}

func testCache(res http.ResponseWriter, req *http.Request) {
	cacheSum := req.URL.Query().Get("sum")
	if cacheSum == "" {
		res.WriteHeader(400)
		_, _ = res.Write([]byte("Please provide ?sum="))
		return
	}

	data := getJsonMap()
	stat := metadata.GetStat(cacheSum)
	if stat == nil {
		data["cached"] = false
		sendJson(data, 200, res)
		return
	}

	data["cached"] = true

	expired, expiresIn := metadata.GetExpired(cacheSum)
	data["expired"] = expired
	data["expires_in_seconds"] = expiresIn

	mimeType := metadata.GetMimeType(cacheSum)
	data["type"] = mimeType

	visitCount := metadata.GetVisits(cacheSum)
	data["visit_count"] = visitCount

	score, size := eviction.ScoreFile(cacheSum)
	data["score"] = score
	data["size_bytes"] = size
	data["size_human"] = humanize.Bytes(uint64(size))

	sendJson(data, 200, res)
}
