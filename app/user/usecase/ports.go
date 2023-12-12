package usecase

import (
	domain2 "user-management/app/user/domain"
)

type InputPorts interface {
	SignUp(request *domain2.SignUpRequest) (*domain2.SignUpResponse, error)
	Login(request *domain2.LoginRequest) (*domain2.LoginResponse, error)
	RefreshToken(request *domain2.RefreshTokenRequest) (*domain2.RefreshTokenResponse, error)
	ChangePassword(request *domain2.ChangePasswordRequest) error
	GetMenuByUserId(request *domain2.MenuRequest) (*domain2.MenuResponse, error)
	GetAllCourseCategory() ([]*domain2.CourseCategoryDatabase, error)
	InsertAuditTrail(request *domain2.AuditTrail) error
}
