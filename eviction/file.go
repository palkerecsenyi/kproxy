package eviction

import (
	"kproxy/metadata"
	"syscall"
	"time"
)

func ScoreFile(fileName string) (float64, int) {
	popularity := metadata.GetVisits(fileName)
	fileInfo := metadata.GetStat(fileName)
	if fileInfo == nil {
		return 0, 0
	}

	size := int(fileInfo.Size())

	var sysInfo syscall.Stat_t
	if sysInfo, ok := fileInfo.Sys().(*syscall.Stat_t); !ok || sysInfo == nil {
		panic("File system doesn't support time metadata.")
	}

	birthTime := time.Unix(sysInfo.Ctimespec.Sec, 0)
	if birthTime == time.Unix(0, 0) {
		panic("File system doesn't support ctime metadata.")
	}

	ageInSeconds := time.Now().Sub(birthTime).Seconds()

	return scoreData(popularity, size, ageInSeconds), size
}
