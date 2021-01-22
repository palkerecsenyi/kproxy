package config

import (
	"github.com/dustin/go-humanize"
	"github.com/prologic/bitcask"
	"io"
	"kproxy/helpers"
	"kproxy/metadata"
	"net/http"
	"os"
)

type DownloadStatus struct {
	Total    uint64
	Database *bitcask.Bitcask
	FileSum  string
}

func (status *DownloadStatus) Write(data []byte) (int, error) {
	bytes := len(data)
	status.Total += uint64(bytes)
	_ = status.Database.Put(
		[]byte(status.FileSum+"-download"),
		[]byte(humanize.Bytes(status.Total)+" downloaded"),
	)
	return bytes, nil
}

func parseRequest(req *http.Request) (string, string) {
	resource := req.URL.Query().Get("url")
	if resource == "" {
		return "", ""
	}

	resourceUrlSum := helpers.UrlSumFromString(resource)
	return resource, resourceUrlSum
}

func downloadLargeFile(res http.ResponseWriter, req *http.Request) {
	resource, resourceUrlSum := parseRequest(req)
	if resource == "" || resourceUrlSum == "" {
		res.WriteHeader(400)
		_, _ = res.Write([]byte("Bad request"))
		return
	}

	resourcePath := helpers.GetObjectPath(resourceUrlSum)
	_ = os.Remove(resourcePath)

	_, _ = res.Write([]byte(
		"Download started! Visit /download-status?url=" + resource,
	))

	db := metadata.GetDatabaseSingleton()
	go func() {
		_ = db.Put([]byte(resourceUrlSum+"-download"), []byte("Download started"))
		response, err := http.Get(resource)
		if err != nil {
			_ = db.Put([]byte(resourceUrlSum+"-download"), []byte(err.Error()))
			return
		}

		output, err := os.Create(resourcePath)
		if err != nil {
			_ = db.Put([]byte(resourceUrlSum+"-download"), []byte("Error saving!"))
			return
		}

		defer func() {
			_ = output.Close()
			response.Body.Close()
		}()

		statusTracker := &DownloadStatus{
			Database: db,
			FileSum:  resourceUrlSum,
		}
		_, err = io.Copy(output, io.TeeReader(response.Body, statusTracker))
		if err != nil {
			_ = db.Put([]byte(resourceUrlSum+"-download"), []byte("Error initialising tracker stream!"))
			return
		}

		_ = db.Put([]byte(resourceUrlSum+"-download"), []byte("Download complete"))
	}()
}

func downloadStatus(res http.ResponseWriter, req *http.Request) {
	resource, resourceUrlSum := parseRequest(req)
	if resource == "" || resourceUrlSum == "" {
		res.WriteHeader(400)
		_, _ = res.Write([]byte("Bad request"))
		return
	}

	db := metadata.GetDatabaseSingleton()
	status, err := db.Get([]byte(resourceUrlSum + "-download"))
	if err != nil {
		res.WriteHeader(404)
		_, _ = res.Write([]byte("Download not found"))
		return
	}

	_, _ = res.Write(status)
}

func downloadSavedFile(res http.ResponseWriter, req *http.Request) {
	resource, resourceUrlSum := parseRequest(req)
	if resource == "" || resourceUrlSum == "" {
		res.WriteHeader(400)
		_, _ = res.Write([]byte("Bad request"))
		return
	}

	file, err := os.Open(helpers.GetObjectPath(resourceUrlSum))
	if err != nil {
		res.WriteHeader(500)
		_, _ = res.Write([]byte("Error opening file! Maybe it doesn't exist."))
		return
	}

	_, err = io.Copy(res, file)
	if err != nil {
		res.WriteHeader(500)
		_, _ = res.Write([]byte("Error starting write stream! Maybe file is corrupt?"))
	}
}
