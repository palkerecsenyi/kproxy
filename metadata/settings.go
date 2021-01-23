package metadata

import (
	"encoding/hex"
	"encoding/json"
	"hash/adler32"
	"net/http"
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

// instead of using req.RemoteAddr, we generate fingerprints using some more reliable data
func GetUserId(req *http.Request) string {
	// OS-specific and browser-specific
	accept := req.Header.Get("Accept")
	// OS-specific, and browser-specific
	userAgent := req.UserAgent()
	// Locale-specific
	acceptLanguage := req.Header.Get("Accept-Language")

	fingerprintData := accept + userAgent + acceptLanguage

	adler := adler32.New()
	_, _ = adler.Write([]byte(fingerprintData))
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
