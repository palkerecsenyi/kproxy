package config

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/dustin/go-humanize"
	"kproxy/certificate"
	"kproxy/eviction"
	"kproxy/helpers"
	"kproxy/metadata"
	"kproxy/metadata/analytics"
	"math"
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
		sendJson(data, 500, false, res)
		return
	}

	cpuUsage := getJsonMap()
	cpuUsage["system"] = stat.CPUStatAll.System
	cpuUsage["user"] = stat.CPUStatAll.User
	cpuUsage["nice"] = stat.CPUStatAll.Nice
	data["cpu"] = cpuUsage

	storageUsage := eviction.CalculateStorageUsage()
	usage := getJsonMap()
	usage["bytes"] = storageUsage
	usage["human"] = humanize.Bytes(uint64(storageUsage))
	data["cache_usage"] = usage

	files, _ := os.ReadDir(helpers.GetPath())
	data["cache_object_count"] = len(files)

	data["your_ip"] = req.RemoteAddr
	data["my_time"] = time.Now().Format(time.RFC3339)

	sendJson(data, 200, false, res)
}

func getLogs(res http.ResponseWriter, req *http.Request) {
	data := getJsonMap()

	daysString := req.URL.Query().Get("days")
	days := 1
	if daysString != "" {
		days, _ = strconv.Atoi(daysString)
	}

	requireCached := req.URL.Query().Get("only-cached")
	logs, lastModified := analytics.GetLogs(time.Now().AddDate(0, 0, 0-days), requireCached == "1")
	formattedLastModified := lastModified.Format(time.RFC3339)
	if ifModifiedSince := parseLastModified(req.Header.Get("if-modified-since")); !ifModifiedSince.IsZero() {

		// only give full response if lastModified is after isModifiedSince
		if !lastModified.After(ifModifiedSince) || ifModifiedSince.Format(time.RFC3339) == formattedLastModified {
			res.WriteHeader(304)
			_, _ = res.Write([]byte("Not Modified"))
			return
		}
	}
	res.Header().Set("last-modified", formattedLastModified)

	totalSavings := analytics.SumSavings(logs)
	fractionCached := analytics.FractionCached(logs)

	data["logs"] = logs
	data["fraction_cached"] = math.Floor(fractionCached*1000) / 1000

	savings := getJsonMap()
	savings["bytes"] = totalSavings
	savings["human"] = humanize.Bytes(totalSavings)
	data["cache_savings"] = savings

	sendJson(data, 200, true, res)
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
		sendJson(data, 200, false, res)
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

	sizeData := getJsonMap()
	sizeData["bytes"] = size
	sizeData["human"] = humanize.Bytes(uint64(size))
	data["size"] = sizeData

	sendJson(data, 200, false, res)
}
