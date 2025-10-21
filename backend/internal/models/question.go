package models

import "time"

type Question struct {
	ID        int       `json:"id" 			db:"id"`
	Title     string    `json:"title" 		db:"title"`
	Answer    string    `json:"answer" 		db:"answer"`
	Tags      []string  `json:"tags" 		db:"tags"`
	CreatedAt time.Time `json:"createdAt" 	db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" 	db:"updatedAt"`
}
