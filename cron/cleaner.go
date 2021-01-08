package cron

import (
	"fmt"
	"kproxy/eviction"
	"kproxy/helpers"
	"kproxy/metadata"
	"os"
	"path"
	"sort"
	"strconv"
)

type fileWithScore struct {
	name  string
	score float64
	size  int64
}

func _stats(removalCount int) {
	fmt.Println("Removed " + strconv.Itoa(removalCount) + " files")
}

func Clean() {
	contents, err := os.ReadDir(helpers.GetPath())
	if err != nil {
		panic(err)
	}

	currentUsage := eviction.CalculateStorageUsage()
	maxUsage := eviction.GetMaxUsage()
	shouldEvict := currentUsage > maxUsage

	removalCount := 0

	var fileScores []fileWithScore
	for _, file := range contents {
		fileName := file.Name()
		filePath := path.Join(helpers.GetPath(), fileName)

		if expired, _ := metadata.GetExpired(fileName); expired {
			_ = os.Remove(filePath)
			removalCount++
			continue
		}

		score, size := eviction.ScoreFile(fileName)
		if score == 0 {
			_ = os.Remove(filePath)
			removalCount++
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
		_stats(removalCount)
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
		removalCount++
		_ = os.Remove(filePath)

		// subtract the usage
		currentUsage -= file.size
		// remove it from the queue
		fileScores = fileScores[1:]
	}

	fmt.Println("Done")
	_stats(removalCount)
}
