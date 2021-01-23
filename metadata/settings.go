package metadata

import (
	"encoding/hex"
	"encoding/json"
	"hash/adler32"
	"net/http"
	"net/url"
)

type CacheRule struct {
	Glob      string
	OnlyTypes []string
	Rule      string // only populated for certain view renders
}

type Settings struct {
	ID          string
	AlwaysCache []CacheRule
	NeverCache  []CacheRule
}

func GetUserId(req *http.Request) string {
	parsedUrl, err := url.Parse("scheme://" + req.RemoteAddr)
	if err != nil {
		parsedUrl, _ = url.Parse("scheme://localhost")
	}

	adler := adler32.New()
	_, _ = adler.Write([]byte(parsedUrl.Hostname()))
	return hex.EncodeToString(adler.Sum(nil))
}

func GetSettings(req *http.Request) Settings {
	settings := Settings{}

	userId := GetUserId(req)
	if userId == "" {
		return settings
	} else {
		settings.ID = userId
	}

	db := GetDatabaseSingleton()
	rawData, err := db.Get([]byte("settings-" + userId))
	if err != nil {
		return settings
	}

	_ = json.Unmarshal(rawData, &settings)
	return settings
}

func (settings *Settings) Save() {
	rawData, err := json.Marshal(settings)
	if err != nil {
		return
	}

	db := GetDatabaseSingleton()
	_ = db.Put([]byte("settings-"+settings.ID), rawData)
}
