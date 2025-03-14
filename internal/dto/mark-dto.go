package dto

type MarkRequestTo struct {
	Name string `json:"name" validate:"required,min=2,max=32"`
}

type MarkResponseTo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
