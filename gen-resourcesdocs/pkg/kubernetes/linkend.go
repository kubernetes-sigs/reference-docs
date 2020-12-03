package kubernetes

import (
	"fmt"
)

// LinkEnds maps definition key to a link-end
type LinkEnds map[Key][]string

// Add a new map between key and linkend
func (o LinkEnds) Add(key Key, linkend []string) {
	o[key] = linkend
}

func (o LinkEnds) Debug() {
	for k, v := range o {
		fmt.Printf("%s: %v\n", k, v)
	}
}
