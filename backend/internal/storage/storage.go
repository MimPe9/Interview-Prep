package storage

import "interviewPrep/backend/internal/models"

type Storage interface {
	GetAllQuestions() ([]models.Question, error)
	CreateQuestion() (*models.Question, error)
	Close() error
	DeleteQuestion(id int) error
}
