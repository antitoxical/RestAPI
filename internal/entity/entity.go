package entity

import (
	"time"
)

type Writer struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"-"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type News struct {
	ID       int64     `json:"id"`
	WriterID int64     `json:"writerId"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

type Message struct {
	ID      int64  `json:"id"`
	NewsID  int64  `json:"newsId"`
	Content string `json:"content"`
}

type Mark struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
