package adapters

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-management/app/domain/course"
)

type ParentCourseModel struct {
	CourseId       primitive.ObjectID `json:"_id" bson:"_id"`
	IsPublished    bool               `json:"is_published"`
	TotalSubCourse int                `json:"total_sub_course,omitempty"`
	Category       *course.Category   `json:"category"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	SubCourses     []*course.Course   `json:"sub_courses"`
}
