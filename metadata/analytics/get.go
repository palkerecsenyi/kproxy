package analytics

import (
	"encoding/json"
	"kproxy/metadata"
	"time"
)

// (unordered requests, last modified)
func GetLogs(since time.Time, requireCached bool) ([]RequestLog, time.Time) {
	var logs []RequestLog
	db := metadata.GetDatabaseSingleton()
	var lastModified time.Time
	_ = db.Scan([]byte("log-"), func(key []byte) error {
		rawData, err := db.Get(key)
		if err != nil {
			return err
		}

		data := RequestLog{}
		_ = json.Unmarshal(rawData, &data)

		if data.Timestamp.Before(time.Now().AddDate(0, 0, -7)) {
			_ = db.Delete(key)
			return nil
		}

		if data.Timestamp.Before(since) {
			return nil
		}

		if requireCached && !data.Cached {
			return nil
		}

		logs = append(logs, data)
		if data.Timestamp.After(lastModified) {
			lastModified = data.Timestamp
		}

		return nil
	})

	return logs, lastModified
}

func SumSavings(logs []RequestLog) uint64 {
	var savings uint64
	for _, log := range logs {
		savings += log.Savings
	}
	return savings
}

func FractionCached(logs []RequestLog) float64 {
	var cached float64
	var uncached float64
	for _, log := range logs {
		if log.Cached {
			cached++
		} else {
			uncached++
		}
	}

	return (cached) / (cached + uncached)
}
