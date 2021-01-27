package config

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func sendJson(data map[string]interface{}, status int, cacheable bool, res http.ResponseWriter) {
	encodedData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	if !cacheable {
		res.Header().Set("Cache-Control", "no-store")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Content-Length", strconv.Itoa(len(encodedData)))
	res.WriteHeader(status)
	_, _ = res.Write(encodedData)
}

func getJsonMap() map[string]interface{} {
	return make(map[string]interface{})
}

func parseLastModified(value string) time.Time {
	attempts := []string{
		time.RFC3339,
		time.RFC822,
		time.RFC850,
		time.RFC1123,
		time.ANSIC,
		time.RFC822Z,
	}

	for _, format := range attempts {
		parsed, err := time.Parse(format, value)
		if err != nil {
			continue
		} else {
			return parsed
		}
	}

	return time.Time{}
}
