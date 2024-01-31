package course

import "sort"

type Course struct {
	Order       int        `json:"order"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Contents    []*Content `json:"contents,omitempty"`
}

func (c *Course) init() {
	sort.SliceStable(c.Contents, func(i, j int) bool {
		return c.Contents[i].Order < c.Contents[j].Order
	})
}
