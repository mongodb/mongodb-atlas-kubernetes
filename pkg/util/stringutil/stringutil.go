package stringutil

import (
	"regexp"
	"sort"
	"strings"
)

// Contains returns true if there is at least one string in `slice`
// that is equal to `s`.
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// SimplifyAndSort removes all characters which are not letters or numbers and sorts the rest
func SimplifyAndSort(str string) string {
	reg := regexp.MustCompile("[^a-z0-9]+")
	temp := reg.ReplaceAllString(strings.ToLower(str), "")
	tempSlice := strings.Split(temp, "")
	sort.Strings(tempSlice)
	return strings.Join(tempSlice, "")
}
