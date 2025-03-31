package dto

type WriterRequestTo struct {
	ID        int64  `json:"id" `
	Login     string `json:"login" validate:"required,min=2,max=64"`
	Password  string `json:"password" validate:"required,min=8,max=128"`
	FirstName string `json:"firstname" validate:"required,min=2,max=64"`
	LastName  string `json:"lastname" validate:"required,min=2,max=64"`
}

type WriterResponseTo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
