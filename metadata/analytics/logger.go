package analytics

import (
	"encoding/json"
	"github.com/google/uuid"
	"kproxy/metadata"
	"net/url"
	"time"
)

type RequestLog struct {
	Cached    bool
	Savings   uint64
	Hostname  string
	Timestamp time.Time
}

func LogRequest(url *url.URL, cached bool, size uint64) {
	log := RequestLog{
		Hostname:  url.Hostname(),
		Timestamp: time.Now(),
	}

	defer func() {
		id := uuid.NewString()
		db := metadata.GetDatabaseSingleton()
		jsonValue, _ := json.Marshal(log)
		_ = db.Put([]byte("log-"+id), jsonValue)
	}()

	if !cached {
		log.Cached = false
		return
	}
	log.Cached = true

	log.Savings = size
}
