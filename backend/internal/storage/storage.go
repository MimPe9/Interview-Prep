package storage

import "interview-prep/backend/internal/models"

type Storage interface {
	GetAllQuestions() ([]models.Question, error)
	CreateQuestion() (*models.Question, error)
	Close() error
}
