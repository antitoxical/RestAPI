package entity

import "time"

type Writer struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Login     string `gorm:"column:login;size:64;not null;unique" json:"login"`
	Password  string `gorm:"column:password;size:128;not null" json:"-"`
	FirstName string `gorm:"column:firstname;size:64;not null" json:"firstname"`
	LastName  string `gorm:"column:lastname;size:64;not null" json:"lastname"`
}

func (News) TableName() string {
	return "tbl_news"
}

func (Mark) TableName() string {
	return "tbl_mark"
}

func (Writer) TableName() string {
	return "tbl_writer"
}

func (Message) TableName() string {
	return "tbl_message"
}

type News struct {
	ID       int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WriterID int64     `gorm:"not null" json:"writerId"`
	Title    string    `gorm:"size:255;not null" json:"title"`
	Content  string    `gorm:"type:text;not null" json:"content"`
	Created  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created"`
	Modified time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"modified"`
	Marks    []Mark    `gorm:"many2many:news_mark;"`
}

type Message struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	NewsID  int64  `gorm:"not null" json:"newsId"`
	Content string `gorm:"type:text;not null" json:"content"`
}

type Mark struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:32;not null;unique" json:"name"`
}
