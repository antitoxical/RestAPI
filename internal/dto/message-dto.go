package dto

type MessageRequestTo struct {
	NewsID  int64  `json:"newsId" validate:"required"`
	Content string `json:"content" validate:"required,min=2,max=2048"`
}

type MessageResponseTo struct {
	ID      int64  `json:"id"`
	NewsID  int64  `json:"newsId"`
	Content string `json:"content"`
}
