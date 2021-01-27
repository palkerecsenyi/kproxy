package config

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func sendJson(data map[string]interface{}, status int, res http.ResponseWriter) {
	encodedData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Cache-Control", "no-store")
	res.Header().Set("Content-Length", strconv.Itoa(len(encodedData)))
	res.WriteHeader(status)
	_, _ = res.Write(encodedData)
}

func getJsonMap() map[string]interface{} {
	return make(map[string]interface{})
}
