package config

import (
	"github.com/dustin/go-humanize"
	"io"
	"kproxy/helpers"
	"kproxy/metadata"
	"net/http"
	"os"
)

type DownloadStatus struct {
	Total   uint64
	FileSum string
}

func (status *DownloadStatus) Write(data []byte) (int, error) {
	bytes := len(data)
	status.Total += uint64(bytes)

	resource := metadata.Get(status.FileSum)
	resource.UpdateDownload(humanize.Bytes(status.Total) + " downloaded")

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

	resourceData := metadata.Get(resourceUrlSum)
	go func() {
		resourceData.UpdateDownload("Download started")

		response, err := http.Get(resource)
		if err != nil {
			resourceData.UpdateDownload(err.Error())
			return
		}

		output, err := os.Create(resourcePath)
		if err != nil {
			resourceData.UpdateDownload("Error saving!")
			return
		}

		defer func() {
			_ = output.Close()
			response.Body.Close()
		}()

		statusTracker := &DownloadStatus{
			FileSum: resourceUrlSum,
		}
		_, err = io.Copy(output, io.TeeReader(response.Body, statusTracker))
		if err != nil {
			resourceData.UpdateDownload("Couldn't initialise status tracker")
			return
		}

		resourceData.UpdateDownload("Download complete!")
	}()
}

func downloadStatus(res http.ResponseWriter, req *http.Request) {
	resource, resourceUrlSum := parseRequest(req)
	if resource == "" || resourceUrlSum == "" {
		res.WriteHeader(400)
		_, _ = res.Write([]byte("Bad request"))
		return
	}

	resourceData := metadata.Get(resourceUrlSum)
	if resourceData.DownloadStatus == "" {
		res.WriteHeader(404)
		_, _ = res.Write([]byte("Download not found"))
		return
	}

	_, _ = res.Write([]byte(resourceData.DownloadStatus))
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
