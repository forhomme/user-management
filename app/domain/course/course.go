package course

import "sort"

type Course struct {
	Order       int
	Title       string
	Description string
	Contents    []*Content
}

func (c *Course) init() {
	sort.SliceStable(c.Contents, func(i, j int) bool {
		return c.Contents[i].Order < c.Contents[j].Order
	})
}
