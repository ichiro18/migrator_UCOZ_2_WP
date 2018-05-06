package helpers

import "regexp"

type StringHelper struct {
}

// Split the given string by a regular expression.
func RegSplit(pattern string, subject string) []string{
	reg := regexp.MustCompile(pattern)
	indexes := reg.FindAllStringIndex(subject, -1)
	laststart := 0
	result := make([]string, len(indexes) + 1)
	for i, element := range indexes {
		result[i] = subject[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = subject[laststart:len(subject)]
	return result
}