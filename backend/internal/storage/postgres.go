package storage

import (
	"context"
	"database/sql"
	"fmt"
	"interview-prep/backend/internal/models"
	"log"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// Подключаемся к БД
func NewPosgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	//Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateQuestion(q *models.Question) error {
	query := `
		INSERT INTO questions (title, answer, tag) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at, updated_at
	`

	return s.db.QueryRow(query, q.Title, q.Answer, q.Tag).Scan(
		&q.ID, &q.CreatedAt, &q.UpdatedAt,
	)
}
func (s *PostgresStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *PostgresStorage) Init(ctx context.Context) error {
	q := `
		CREATE TABLE IF NOT EXISTS interview_prep (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
		answer TEXT NOT NULL,
		tags TEXT[],
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
	`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}

func (s *PostgresStorage) GetAllQuestions() ([]models.Question, error) {
	q := `
		SELECT id, title, answer, tag, created_at, updated_at
        FROM questions
        ORDER BY tag DESC
	`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var question []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(&q.ID, &q.Title, &q.Answer, &q.Tag, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		question = append(question, q)
	}

	return question, nil
}
