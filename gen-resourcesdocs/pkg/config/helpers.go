package config

import "strings"

func getEscapedFirstPhrase(s string) string {
	description := strings.Split(s, ".")[0] + "."
	return strings.ReplaceAll(description, "\"", "\\\"")
}
