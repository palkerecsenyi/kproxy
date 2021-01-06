package cache

import (
	"bytes"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Get(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	if !shouldGetFromCache(req) {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	urlSum := getUrlSum(ctx)
	data, err := os.ReadFile(getObjectPath(urlSum))
	if err != nil {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	response := &http.Response{}
	response.Header = make(http.Header)
	response.Header.Add("X-Cache-Sum", urlSum)
	response.Header.Set("Cache-Control", "no-cache")

	contentType := http.DetectContentType(data)
	response.Header.Set("Content-Type", contentType)
	var dataBuffer *bytes.Buffer
	if strings.HasPrefix(contentType, "text/") {
		dataBuffer = bytes.NewBufferString(string(data))
	} else {
		dataBuffer = bytes.NewBuffer(data)
	}

	response.ContentLength = int64(dataBuffer.Len())
	response.Body = ioutil.NopCloser(dataBuffer)

	response.Request = req
	response.TransferEncoding = req.TransferEncoding
	response.StatusCode = http.StatusOK
	response.Status = http.StatusText(http.StatusOK)

	ctx.UserData = ProxyCacheState{
		FromCache: true,
	}

	return response
}
