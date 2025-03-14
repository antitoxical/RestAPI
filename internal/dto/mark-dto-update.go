package dto

type MarkUpdateRequestTo struct {
	Name string `json:"name" validate:"required,min=2,max=32"`
	ID   int64  `json:"id"`
}
