package infrastructure

import (
	domain2 "user-management/app/user/domain"
)

type DatabasePorts interface {
	InsertAuditTrail(trail *domain2.AuditTrail) error
	GetUserById(id string) (*domain2.UsersDatabase, error)
	GetUserByEmail(email string) (*domain2.UsersDatabase, error)
	InsertUser(*domain2.UsersDatabase) (*domain2.UsersDatabase, error)
	UpdateUserById(id string, dataUser *domain2.UsersDatabase) error
	GetUserMenu(roleId int) ([]*domain2.Menu, error)
	GetAllCourseCategory() ([]*domain2.CourseCategoryDatabase, error)
	InsertCourseContent(in *domain2.CourseDatabase) ([]*domain2.CourseDatabase, error)
}
