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
	if !shouldGetFromCache(req) {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	urlSum := helpers.GetUrlSum(ctx)
	metadata.IncrementVisits(urlSum)

	contentType := metadata.GetMimeType(urlSum)
	// avoid unexpected behaviour by not assuming mime types
	if contentType == "" {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	expired, expiresInSeconds := metadata.GetExpired(urlSum)
	if expired {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	data, err := os.ReadFile(helpers.GetObjectPath(urlSum))
	if err != nil {
		ctx.UserData = ProxyCacheState{
			FromCache: false,
		}
		return nil
	}

	response := &http.Response{}
	response.Header = make(http.Header)
	response.Header.Add("X-Cache-Sum", urlSum)
	response.Header.Add("X-Cache-Expires-In", strconv.Itoa(expiresInSeconds))
	// allow browser caching for burst periods of one hour to reduce proxy load
	response.Header.Set("Cache-Control", "max-age=3600")
	// since the server didn't specify the request as being private, we can tell the browser it's public, too
	response.Header.Add("Cache-Control", "public")

	// yes this is horrifyingly unsafe
	if origin := ctx.Req.Header.Get("Origin"); origin != "" {
		response.Header.Add("Access-Control-Allow-Origin", origin)
		response.Header.Add("Access-Control-Allow-Credentials", "true")
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

	ctx.UserData = ProxyCacheState{
		FromCache: true,
	}

	return response
}
