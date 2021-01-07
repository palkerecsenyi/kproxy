package cron

import (
	"fmt"
	"kproxy/eviction"
	"kproxy/helpers"
	"kproxy/metadata"
	"os"
	"path"
	"sort"
)

type fileWithScore struct {
	name  string
	score float64
	size  int64
}

func Clean() {
	contents, err := os.ReadDir(helpers.GetPath())
	if err != nil {
		panic(err)
	}

	currentUsage := eviction.CalculateStorageUsage()
	maxUsage := eviction.GetMaxUsage()
	shouldEvict := currentUsage > maxUsage

	var fileScores []fileWithScore
	for _, file := range contents {
		fileName := file.Name()
		filePath := path.Join(helpers.GetPath(), fileName)

		if metadata.GetExpired(fileName) {
			_ = os.Remove(filePath)
			continue
		}

		score, size := eviction.ScoreFile(fileName)
		if score == 0 {
			_ = os.Remove(filePath)
			continue
		}

		if shouldEvict {
			fileScores = append(fileScores, fileWithScore{
				name:  fileName,
				score: score,
				size:  int64(size),
			})
		}
	}

	if !shouldEvict {
		return
	}

	fmt.Println("Evicting files until KPROXY_MAX_SPACE is reached")

	// sort in ascending order of scores
	sort.Slice(fileScores, func(i, j int) bool {
		return fileScores[i].score < fileScores[j].score
	})

	// while we're over our usage limit
	for currentUsage > maxUsage && len(fileScores) > 0 {
		// delete the smallest-scored file from system
		file := fileScores[0]
		filePath := path.Join(helpers.GetPath(), file.name)
		_ = os.Remove(filePath)

		// subtract the usage
		currentUsage -= file.size
		// remove it from the queue
		fileScores = fileScores[1:]
	}

	fmt.Println("Done")
}
