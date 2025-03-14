package dto

type MessageUpdateRequestTo struct {
	NewsID  int64  `json:"newsId" validate:"required"`
	Content string `json:"content" validate:"required,min=2,max=2048"`
	ID      int64  `json:"id"`
}
