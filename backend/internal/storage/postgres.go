package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"interviewPrep/backend/internal/models"

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
	log.Printf("Creating question: %s", q.Title)

	var exists bool

	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE title = $1)", q.Title).Scan(&exists)
	if err != nil {
		log.Printf("Error checking question existence: %v", err)
		return fmt.Errorf("error checking question existence: %w", err)
	}

	if exists {
		log.Printf("Question with title already exists: %s", q.Title)
		return fmt.Errorf("question with this title already exists")
	}

	query := `
		INSERT INTO questions (title, answer, tags) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at, updated_at
	`

	err = s.db.QueryRow(query, q.Title, q.Answer, pq.Array(q.Tags)).Scan(
		&q.ID, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating question: %v", err)
		return fmt.Errorf("error creating question: %w", err)
	}

	log.Printf("Question created successfully with ID: %d", q.ID)
	return nil
}

func (s *PostgresStorage) UpdateQuestion(q *models.Question) error {
	log.Printf("Updating question ID: %d", q.ID)

	query := `
		UPDATE questions
		SET title = $1, answer = $2, tags = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	res, err := s.db.Exec(query, q.Title, q.Answer, pq.Array(q.Tags), q.ID)
	if err != nil {
		log.Printf("Error updating question ID %d: %v", q.ID, err)
		return fmt.Errorf("failed to update: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for question ID %d: %v", q.ID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("Question not found for update ID: %d", q.ID)
		return fmt.Errorf("question with id %d not found", q.ID)
	}

	log.Printf("Question updated successfully ID: %d", q.ID)
	return nil
}

func (s *PostgresStorage) DeleteQuestion(id int) error {
	log.Printf("Deleting question ID: %d", id)

	query := `DELETE FROM questions WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting question ID %d: %v", id, err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for deletion ID %d: %v", id, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Question not found for deletion ID: %d", id)
		return fmt.Errorf("question with id %d not found", id)
	}

	log.Printf("Question deleted successfully ID: %d", id)
	return nil
}

func (s *PostgresStorage) Close() error {
	log.Println("Closing database connection")
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *PostgresStorage) GetAllQuestions() ([]models.Question, error) {
	log.Println("Getting all questions")

	q := `
		SELECT id, title, answer, tags, created_at, updated_at
        FROM questions
        ORDER BY created_at DESC
	`

	rows, err := s.db.Query(q)
	if err != nil {
		log.Printf("Error querying all questions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		var tags string
		err := rows.Scan(&q.ID, &q.Title, &q.Answer, &tags, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning question row: %v", err)
			return nil, err
		}

		q.Tags = parse(tags)
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over question rows: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d questions", len(questions))
	return questions, nil
}

func (s *PostgresStorage) GetQuestionByID(id int) (*models.Question, error) {
	log.Printf("Getting question by ID: %d", id)

	q := `
		SELECT id, title, answer, tags, created_at, updated_at 
		FROM questions
		WHERE id = $1
	`

	var question models.Question
	var tags string
	err := s.db.QueryRow(q, id).Scan(&question.ID, &question.Title, &question.Answer, &tags, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Question not found ID: %d", id)
			return nil, fmt.Errorf("question not found")
		}
		log.Printf("Error getting question ID %d: %v", id, err)
		return nil, err
	}

	question.Tags = parse(tags)
	log.Printf("Question found ID: %d", id)
	return &question, nil
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
