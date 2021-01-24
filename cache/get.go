package cache

import (
	"bytes"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"kproxy/helpers"
	"kproxy/metadata"
	"net/http"
	"os"
	"strconv"
)

func Get(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	userData := ProxyCacheState{
		FromCache: false,
	}

	urlSum := metadata.ClientUrlSum(req.URL.String(), req.Header)
	contentType := metadata.GetMimeType(urlSum)
	// avoid unexpected behaviour by not assuming mime types
	if contentType == "" {
		ctx.UserData = userData
		return nil
	}

	if !shouldGetFromCache(req, contentType, urlSum) {
		ctx.UserData = userData
		return nil
	}

	metadata.IncrementVisits(urlSum)

	expired, expiresInSeconds := metadata.GetExpired(urlSum)
	if expired {
		ctx.UserData = userData
		return nil
	}

	data, err := os.ReadFile(helpers.GetObjectPath(urlSum))
	if err != nil {
		ctx.UserData = userData
		return nil
	}

	response := &http.Response{}
	response.Header = metadata.GetHeaders(urlSum)

	response.Header.Add("X-Cache-Sum", urlSum)
	response.Header.Add("X-Cache-Expires-In", strconv.Itoa(expiresInSeconds))
	// allow browser caching for burst periods of one hour to reduce proxy load
	response.Header.Set("Cache-Control", "max-age=3600")
	// since the server didn't specify the request as being private, we can tell the browser it's public, too
	response.Header.Add("Cache-Control", "public")

	// yes this is horrifyingly unsafe
	if origin := ctx.Req.Header.Get("Origin"); origin != "" {
		response.Header.Add("Access-Control-Allow-Origin", origin)
	}

	response.Header.Set("Content-Type", contentType)

	var dataBuffer *bytes.Buffer
	if helpers.IsTextualMime(contentType) {
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

	userData.FromCache = true
	ctx.UserData = userData

	return response
}
