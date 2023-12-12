package _interface

type Category struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type AllCategory struct {
	Category []*Category `json:"category"`
}

type GetCourses struct {
	Id         string
	CategoryId int    `json:"CategoryId"`
	Filter     string `json:"Filter"`
	PerPage    int    `json:"PerPage"`
	Page       int    `json:"Page"`
}

type AllCourse struct {
	Courses []*Course `json:"Courses"`
}

type Course struct {
	CourseId    string
	CategoryId  int
	Title       string
	Description string
	Tags        []string
	SubCourses  []*SubCourse
}

type SubCourse struct {
	Title       string
	Description string
	Contents    []*Content
}

type Content struct {
	IsAssignment bool   `json:"is_assignment"`
	Ordering     int    `json:"ordering"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Content      string `json:"content"`
}
