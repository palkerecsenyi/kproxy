package cache

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func getPath() string {
	rootPath := os.Getenv("KPROXY_PATH")
	if rootPath == "" {
		panic("KPROXY_PATH is unset")
	}

	return rootPath
}

func getObjectPath(object string) string {
	rootPath := getPath()
	return path.Join(rootPath, object)
}

func responseToBytes(response *http.Response) []byte {
	buf, _ := ioutil.ReadAll(response.Body)
	cacheReader := ioutil.NopCloser(bytes.NewBuffer(buf))
	defer cacheReader.Close()
	responseReader := ioutil.NopCloser(bytes.NewBuffer(buf))
	response.Body = responseReader

	cacheBuffer := new(bytes.Buffer)
	_, _ = cacheBuffer.ReadFrom(cacheReader)

	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/html"
	}

	if strings.HasPrefix(contentType, "text/") || contentType == "application/javascript" {
		return []byte(cacheBuffer.String())
	} else {
		return cacheBuffer.Bytes()
	}
}

func getUrlSum(ctx *goproxy.ProxyCtx) string {
	urlSha1 := sha1.New()
	urlSha1.Write([]byte(ctx.Req.URL.String()))
	return hex.EncodeToString(urlSha1.Sum(nil))
}

func getMaxAge(response *http.Response) time.Duration {
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
			return time.Duration(0)
		}

		maxAgeValue, err := strconv.Atoi(maxAgeSlice[1])
		if err != nil {
			return time.Duration(0)
		}

		return time.Duration(maxAgeValue) * time.Second
	}

	return time.Duration(0)
}
