package dto

type NewsUpdateRequestTo struct {
	WriterID int64                 `json:"writerId" validate:"required"`
	Title    string                `json:"title" validate:"required,min=2,max=64"`
	Content  string                `json:"content" validate:"required,min=4,max=2048"`
	ID       int64                 `json:"id"`
	Marks    []MarkUpdateRequestTo `json:"marks"`
}
