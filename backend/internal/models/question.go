package models

import "time"

type Question struct {
	ID        int
	Title     string
	Answer    string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
