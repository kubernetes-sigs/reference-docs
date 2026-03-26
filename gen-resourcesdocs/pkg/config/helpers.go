package config

import (
	"regexp"
	"strings"
)

// sentenceEnd matches a period that is likely a real sentence boundary:
// it must be followed by a space (or end of string) and NOT preceded by
// a single uppercase letter (which would indicate an abbreviation like X.509)
var sentenceEnd = regexp.MustCompile(`(?:[^A-Z])\.\s`)

func getEscapedFirstPhrase(s string) string {
	loc := sentenceEnd.FindStringIndex(s)
	var description string
	if loc != nil {
		description = s[:loc[0]+2]
	} else {
		description = s
	}
	return strings.ReplaceAll(description, "\"", "\\\"")
}
