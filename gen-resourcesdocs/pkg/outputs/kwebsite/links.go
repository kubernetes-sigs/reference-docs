package kwebsite

import (
	"fmt"
	"strings"
)

// LinkEnd returns a link to a section in a part/chapter
// s is an array containing partname / chaptername
func (o *KWebsite) LinkEnd(s []string, name string) string {
	typename := name
	mapprefix := ""
	array := ""

	if strings.HasPrefix(typename, "map[string]") {
		mapprefix = "map[string]"
		typename = strings.TrimPrefix(name, mapprefix)
	}

	if strings.HasPrefix(typename, "[]") {
		array = "[]"
		typename = strings.TrimPrefix(name, array)
	}

	return fmt.Sprintf("%s%s<a href=\"{{< ref \"../%s/%s#%s\" >}}\">%s</a>", mapprefix, array, escapeName(s[0]), escapeName(s[1]), headingID(typename), typename)
}
