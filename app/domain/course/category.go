package course

type Category struct {
	CategoryId   int    `json:"category_id" bson:"category_id"`
	CategoryName string `json:"category_name" bson:"category_name" validate:"required"`
}
