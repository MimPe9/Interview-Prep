package storage

import (
	"database/sql"
	"fmt"
	"interview-prep/backend/internal/models"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// Подключаемся к БД
func NewPosgresStorage(connStr string) (*PostgresStorage, error) {
	log.Printf("Connection string: %s", connStr)

	var db *sql.DB
	var err error

	// Попытки подключения с задержкой
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: sql.Open error: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return &PostgresStorage{db: db}, nil
		}

		log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
}

func (s *PostgresStorage) CreateQuestion(q *models.Question) error {
	query := `
		INSERT INTO questions (title, answer, tags) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at, updated_at
	`

	return s.db.QueryRow(query, q.Title, q.Answer, pq.Array(q.Tags)).Scan(
		&q.ID, &q.CreatedAt, &q.UpdatedAt,
	)
}
func (s *PostgresStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *PostgresStorage) GetAllQuestions() ([]models.Question, error) {
	q := `
		SELECT id, title, answer, tags, created_at, updated_at
        FROM questions
        ORDER BY tags DESC
	`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var question []models.Question
	for rows.Next() {
		var q models.Question
		var tags string
		err := rows.Scan(&q.ID, &q.Title, &q.Answer, &tags, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}

		q.Tags = parse(tags)
		question = append(question, q)
	}

	return question, nil
}

func parse(str string) []string {
	if str == "" {
		return []string{}
	}

	cleaned := strings.Trim(str, "{}")
	if cleaned == "" {
		return []string{}
	}

	rawTags := strings.Split(cleaned, ",")
	tags := make([]string, 0, len(rawTags))

	for _, tag := range rawTags {
		trimmed := strings.TrimSpace(tag)
		trimmed = strings.Trim(trimmed, `"`)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	return tags
}
