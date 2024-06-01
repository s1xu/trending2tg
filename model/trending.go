package model

import (
	"time"
)

type Trends struct {
	ID       uint      `db:"id" json:"id"`
	Title    string    `db:"title" json:"title"`
	Language string    `db:"language" json:"language"`
	Url      string    `db:"url" json:"url"`
	CreateAt time.Time `db:"create_at" json:"create_at"`
}

func NewTrends(title, language, url string) Trends {
	return Trends{
		Title:    title,
		Language: language,
		Url:      url,
		CreateAt: time.Now(),
	}
}
