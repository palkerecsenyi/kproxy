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

func responseToString(response *http.Response) string {
	buf, _ := ioutil.ReadAll(response.Body)
	cacheReader := ioutil.NopCloser(bytes.NewBuffer(buf))
	responseReader := ioutil.NopCloser(bytes.NewBuffer(buf))
	response.Body = responseReader

	cacheBuffer := new(bytes.Buffer)
	_, _ = cacheBuffer.ReadFrom(cacheReader)

	return cacheBuffer.String()
}

func getUrlSum(ctx *goproxy.ProxyCtx) string {
	urlSha1 := sha1.New()
	urlSha1.Write([]byte(ctx.Req.URL.String()))
	return hex.EncodeToString(urlSha1.Sum(nil))
}
