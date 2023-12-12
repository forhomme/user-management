package model

import (
	"sort"
)

type CoursePath struct {
	IsPublished    bool
	TotalSubCourse int
	Category       *Category
	CourseId       string
	Title          string
	Description    string
	SubCourses     []*Courses
}

type Courses struct {
	Title       string
	Description string
	Contents    []*Content
}

type Content struct {
	IsAssignment bool
	Ordering     int
	Content      string
	Title        string
	Description  string
}

func (c *Courses) init() {
	sort.SliceStable(c.Contents, func(i, j int) bool {
		return c.Contents[i].Ordering < c.Contents[j].Ordering
	})
}

func (c *CoursePath) init() {
	c.TotalSubCourse = len(c.SubCourses)
}

func (c *CoursePath) Publish() {
	c.IsPublished = true
}

func (c *CoursePath) IsCourseVisible() bool {
	return c.IsPublished
}
