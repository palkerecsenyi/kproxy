package helpers

import "strings"

func SliceIterator(iterator func(value string) bool, slice []string) bool {
	for _, value := range slice {
		if iterator(value) {
			return true
		}
	}
	return false
}

func SliceContainsPrefix(searchValue string, slice []string) bool {
	return SliceIterator(func(value string) bool {
		return strings.HasPrefix(searchValue, value)
	}, slice)
}

func SliceContainsString(searchValue string, slice []string) bool {
	return SliceIterator(func(value string) bool {
		return value == searchValue
	}, slice)
}

func SliceContainsAnyString(slice []string, prefix ...string) bool {
	return SliceIterator(func(value string) bool {
		return SliceContainsString(value, slice)
	}, prefix)
}
