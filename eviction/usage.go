package eviction

import (
	"io/fs"
	"kproxy/helpers"
	"os"
	"path/filepath"
	"strconv"
)

func CalculateStorageUsage() int64 {
	storageRoot := helpers.GetPath()
	var size int64 = 0
	_ = filepath.Walk(storageRoot, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// shouldn't happen anyway
		if info.IsDir() {
			return nil
		}

		size += info.Size()
		return nil
	})

	return size
}

func GetMaxUsage() int64 {
	maxUsageString := os.Getenv("KPROXY_MAX_SPACE")
	if maxUsageString == "" {
		panic("KPROXY_MAX_SPACE is not set")
	}

	maxUsage, err := strconv.Atoi(maxUsageString)
	if err != nil {
		panic("Could not parse KPROXY_MAX_SPACE")
	}

	return int64(maxUsage)
}
