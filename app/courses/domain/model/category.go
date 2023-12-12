package model

type Category struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name" validate:"required"`
}
