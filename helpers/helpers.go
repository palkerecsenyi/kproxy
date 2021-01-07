package helpers

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/elazarl/goproxy"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetDatabasePath() string {
	databasePath := os.Getenv("KPROXY_DB_PATH")
	if databasePath == "" {
		panic("KPROXY_DB_PATH is unset")
	}

	return databasePath
}

func GetPath() string {
	rootPath := os.Getenv("KPROXY_PATH")
	if rootPath == "" {
		panic("KPROXY_PATH is unset")
	}

	return rootPath
}

func GetObjectPath(object string) string {
	rootPath := GetPath()
	return path.Join(rootPath, object)
}

func GetMimeTypeFromHeader(response *http.Response) string {
	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/html"
	}
	return contentType
}

func IsTextualMime(mime string) bool {
	if strings.HasPrefix(mime, "text/") {
		return true
	}

	if SliceContainsPrefix(mime, []string{
		"application/javascript",
		"image/svg+xml",
	}) {
		return true
	}

	return false
}

func ResponseToBytes(response *http.Response) []byte {
	buf, _ := io.ReadAll(response.Body)
	cacheReader := io.NopCloser(bytes.NewBuffer(buf))
	defer cacheReader.Close()

	// duplicate the buffer for the actual response itself
	responseReader := io.NopCloser(bytes.NewBuffer(buf))
	response.Body = responseReader

	cacheBuffer := new(bytes.Buffer)
	size, err := cacheBuffer.ReadFrom(cacheReader)

	var maxObjectSizeMegabytes int64 = 100
	if err != nil || size > maxObjectSizeMegabytes*1000000 {
		return nil
	}

	contentType := GetMimeTypeFromHeader(response)

	if IsTextualMime(contentType) {
		return []byte(cacheBuffer.String())
	} else {
		return cacheBuffer.Bytes()
	}
}

func GetUrlSum(ctx *goproxy.ProxyCtx) string {
	urlSha1 := sha1.New()
	urlSha1.Write([]byte(ctx.Req.URL.String()))
	return hex.EncodeToString(urlSha1.Sum(nil))
}

func GetRequestMaxAge(response *http.Response) time.Duration {
	cacheControlHeader := response.Header.Get("Cache-Control")
	if cacheControlHeader == "" {
		return time.Duration(0)
	}

	slice := strings.Split(cacheControlHeader, ", ")
	for _, key := range slice {
		if !strings.HasPrefix(key, "max-age") {
			continue
		}

		maxAgeSlice := strings.Split(key, "=")

		if len(maxAgeSlice) != 2 {
			return time.Duration(-1)
		}

		maxAgeValue, err := strconv.Atoi(maxAgeSlice[1])
		if err != nil {
			return time.Duration(-1)
		}

		return time.Duration(maxAgeValue) * time.Second
	}

	return time.Duration(-1)
}
