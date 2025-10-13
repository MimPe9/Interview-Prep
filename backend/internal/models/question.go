package models

import "time"

type Question struct {
	ID        int
	Title     string
	Answer    string
	Tag       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
