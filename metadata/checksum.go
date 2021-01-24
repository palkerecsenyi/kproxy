package metadata

import (
	"encoding/hex"
	"hash/adler32"
	"kproxy/helpers"
	"net/http"
	"strings"
)

func stringToSum(data string) string {
	adler := adler32.New()
	_, _ = adler.Write([]byte(data))
	return hex.EncodeToString(adler.Sum(nil))
}

func variableUrlSum(url string, clientHeaders, serverHeaders http.Header) string {
	variableHeaders := helpers.DecodeMultivalueHeader(serverHeaders.Values("Vary"))

	var clientVariableHeaders []string
	for _, headerName := range variableHeaders {
		clientHeaderValue := clientHeaders.Get(headerName)
		if clientHeaderValue == "" {
			continue
		}

		clientVariableHeaders = append(clientVariableHeaders, headerName+"="+clientHeaderValue+" ")
	}

	checksumAddon := strings.Join(clientVariableHeaders, "")
	if len(checksumAddon) == 0 {
		return stringToSum(url)
	}

	return stringToSum(url + checksumAddon)
}

func ClientUrlSum(url string, clientHeaders http.Header) string {
	resource := Get(stringToSum(url))
	return variableUrlSum(url, clientHeaders, resource.Headers)
}

func ServerUrlSum(url string, clientHeaders, serverHeaders http.Header) string {
	return variableUrlSum(url, clientHeaders, serverHeaders)
}
