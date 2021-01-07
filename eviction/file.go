package eviction

import (
	"kproxy/metadata"
	"time"
)

func ScoreFile(fileName string) (float64, int) {
	popularity := metadata.GetVisits(fileName)
	fileInfo := metadata.GetStat(fileName)
	if fileInfo == nil {
		return 0, 0
	}

	size := int(fileInfo.Size())

	birthTime := fileInfo.ModTime()
	if birthTime == time.Unix(0, 0) {
		panic("File system doesn't support ctime metadata.")
	}

	ageInSeconds := time.Now().Sub(birthTime).Seconds()

	return scoreData(popularity, size, ageInSeconds), size
}
