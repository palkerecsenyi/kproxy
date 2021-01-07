package eviction

import (
	"kproxy/metadata"
	"syscall"
	"time"
)

func ScoreFile(fileName string) float64 {
	popularity := metadata.GetVisits(fileName)
	fileInfo := metadata.GetStat(fileName)
	if fileInfo == nil {
		return 0
	}

	size := int(fileInfo.Size())

	var sysInfo syscall.Stat_t
	if sysInfo, ok := fileInfo.Sys().(*syscall.Stat_t); !ok || sysInfo == nil {
		panic("File system doesn't support time metadata.")
	}

	birthTime := time.Unix(sysInfo.Birthtimespec.Sec, 0)
	if birthTime == time.Unix(0, 0) {
		panic("File system doesn't support birthtime metadata.")
	}

	ageInSeconds := time.Now().Sub(birthTime).Seconds()

	return scoreData(popularity, size, ageInSeconds)
}
