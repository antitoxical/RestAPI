package dto

import (
	"time"
)

type NewsRequestTo struct {
	WriterID int64    `json:"writerId" validate:"required"`
	Title    string   `json:"title" validate:"required,min=2,max=64"`
	Content  string   `json:"content" validate:"required,min=4,max=2048"`
	Marks    []string `json:"marks"`
}

type NewsResponseTo struct {
	ID       int64            `json:"id"`
	WriterID int64            `json:"writerId"`
	Title    string           `json:"title"`
	Content  string           `json:"content"`
	Created  time.Time        `json:"created"`
	Modified time.Time        `json:"modified"`
	Marks    []MarkResponseTo `json:"marks"`
}
