package storage

import "interviewPrep/backend/internal/models"

type Storage interface {
	GetAllQuestions() ([]models.Question, error)
	CreateQuestion(q *models.Question) error
	Close() error
	DeleteQuestion(id int) error
	GetQuestionByID(id int) (*models.Question, error)
	UpdateQuestion(q *models.Question) error
}
