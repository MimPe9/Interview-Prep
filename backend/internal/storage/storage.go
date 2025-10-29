package storage

import "github.com/MimPe9/interview-prep/backend/internal/models"

type Storage interface {
	GetAllQuestions() ([]models.Question, error)
	CreateQuestion() (*models.Question, error)
	Close() error
	DeleteQuestion(id int) error
}
