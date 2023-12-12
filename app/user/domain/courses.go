package domain

import "mime/multipart"

type CourseCategoryDatabase struct {
	CategoryId   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type InsertCourseRequest struct {
	IsPublish   bool   `json:"is_publish"`
	CategoryId  int    `json:"category_id"`
	Ordering    int    `json:"ordering"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     struct {
		File   multipart.File
		Header *multipart.FileHeader
	} `json:"content"`
	Image struct {
		File   multipart.File
		Header *multipart.FileHeader
	} `json:"image"`
	Video struct {
		File   multipart.File
		Header *multipart.FileHeader
	} `json:"video"`
	Tags []string `json:"tags"`
}

type InsertSubContentRequest struct {
	ParentId string `json:"parent_id"`
	InsertCourseRequest
}

type GetCourseContent struct {
	CategoryId int    `json:"category_id"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Filter     string `json:"filter"`
}

type CourseDatabase struct {
	IsPublish      bool     `json:"is_publish"`
	CategoryId     int      `json:"category_id"`
	Ordering       int      `json:"ordering"`
	TotalSubCourse int      `json:"total_sub_course"`
	Id             string   `json:"id"`
	ParentId       string   `json:"parent_id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Content        string   `json:"content"`
	Image          string   `json:"image"`
	Video          string   `json:"video"`
	Tags           []string `json:"tags"`
	CreatedBy      string   `json:"created_by"`
	UpdatedBy      string   `json:"updated_by"`
}
