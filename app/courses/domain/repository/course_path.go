package repository

type CoursePathRepository interface {
	Publish()
	IsCourseVisible() bool
}
