package course

import "sort"

type CoursePath struct {
	IsPublished    bool      `json:"is_published"`
	TotalSubCourse int       `json:"total_sub_course,omitempty"`
	Category       *Category `json:"category"`
	CourseId       string    `json:"course_id,omitempty"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	SubCourses     []*Course `json:"sub_courses"`
}

func (c *CoursePath) List() {
	c.TotalSubCourse = len(c.SubCourses)
	sort.SliceStable(c.SubCourses, func(i, j int) bool {
		return c.SubCourses[i].Order < c.SubCourses[j].Order
	})
}

func (c *CoursePath) Replace(cp *CoursePath) {
	*c = *cp
}

func (c *CoursePath) Publish() {
	c.IsPublished = true
}

func (c *CoursePath) IsCourseVisible() bool {
	return c.IsPublished
}

type FilterCourse struct {
	ID         string `json:"id"`
	Filter     string `json:"filter"`
	CategoryId int    `json:"category_id"`
	Page       int64  `json:"page"`
	PerPage    int64  `json:"per_page"`
	User       *User
}

func (f *FilterCourse) init() {
	if f.Page == 0 {
		f.Page = 1
	}

	if f.PerPage == 0 {
		f.PerPage = 10
	}
}
