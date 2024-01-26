package adapters

import "go.mongodb.org/mongo-driver/bson/primitive"

type ParentCourseModel struct {
	CourseId    primitive.ObjectID `json:"_id" bson:"_id"`
	CategoryId  string             `json:"category_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Tags        []string           `json:"tags"`
	IsPublish   int                `json:"is_publish"`
	SubCourses  []*CourseModel     `json:"sub_courses"`
}

type CourseModel struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Contents    []*ContentModel `json:"contents"`
}

type ContentModel struct {
	NeedLogged  bool   `json:"need_logged"`
	Ordering    int    `json:"ordering"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}
