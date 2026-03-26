package config

import (
	"regexp"
	"strings"
)

// Matches ". " or ".\n" followed by an uppercase letter — a real sentence boundary.
// Avoids splitting on abbreviations like "X.509", "e.g.", "i.e.".
var sentenceEndRe = regexp.MustCompile(`\.\s+[A-Z]`)

func getEscapedFirstPhrase(s string) string {
	if loc := sentenceEndRe.FindStringIndex(s); loc != nil {
		return escape(s[:loc[0]+1])
	}
	s = strings.TrimRight(s, " \t\n")
	if len(s) > 0 && s[len(s)-1] != '.' {
		s += "."
	}
	return escape(s)
}

func escape(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
