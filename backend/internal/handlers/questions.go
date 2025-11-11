package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"interviewPrep/backend/internal/cache"
	"interviewPrep/backend/internal/models"
	"interviewPrep/backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	storage *storage.PostgresStorage
	cache   *cache.RedisCache
}

func NewQuestionHandler(storage *storage.PostgresStorage, cache *cache.RedisCache) *QuestionHandler {
	return &QuestionHandler{
		storage: storage,
		cache:   cache,
	}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	log.Printf("API: Creating new question")

	var question models.Question
	if err := c.BindJSON(&question); err != nil {
		log.Printf("API ERROR: Invalid JSON input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	log.Printf("API: Question data received - Title: %s, Tags: %v", question.Title, question.Tags)

	if err := h.storage.CreateQuestion(&question); err != nil {
		log.Printf("API ERROR: Failed to create question - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create question"})
		return
	}

	// Инвалидируем кэш списка вопросов
	if h.cache != nil {
		if err := h.cache.Delete("questions:all"); err != nil {
			log.Printf("API WARNING: Failed to invalidate cache - %v", err)
		} else {
			log.Printf("API: Cache invalidated for questions:all")
		}
	}

	log.Printf("API SUCCESS: Question created with ID: %d", question.ID)
	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("API ERROR: Invalid ID parameter - %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	log.Printf("API: Updating question ID: %d", id)

	var question models.Question
	if err := c.BindJSON(&question); err != nil {
		log.Printf("API ERROR: Invalid JSON input for update - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Устанавливаем ID из URL параметра
	question.ID = id

	log.Printf("API: Update data - Title: %s, Tags: %v", question.Title, question.Tags)

	if err := h.storage.UpdateQuestion(&question); err != nil {
		log.Printf("API ERROR: Failed to update question ID %d - %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	// Инвалидируем кэш
	if h.cache != nil {
		h.cache.Delete("questions:all")
		h.cache.Delete(fmt.Sprintf("question:%d", id))
		log.Printf("API: Cache invalidated for questions:all and question:%d", id)
	}

	log.Printf("API SUCCESS: Question updated ID: %d", id)
	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	log.Printf("API: Getting all questions")

	cacheKey := "questions:all"

	// Пытаемся получить из кэша
	if h.cache != nil {
		var questions []models.Question
		err := h.cache.Get(cacheKey, &questions)
		if err == nil {
			log.Printf("API: Returning %d questions from cache", len(questions))
			c.JSON(http.StatusOK, questions)
			return
		}
		log.Printf("API: Cache miss for %s - %v", cacheKey, err)
	}

	// Если нет в кэше, то из БД
	questions, err := h.storage.GetAllQuestions()
	if err != nil {
		log.Printf("API ERROR: Failed to get questions from storage - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Questions not found"})
		return
	}

	// Оставляем в кэше на 10 минут
	if h.cache != nil {
		if err := h.cache.Set(cacheKey, questions, 10*time.Minute); err != nil {
			log.Printf("API WARNING: Failed to set cache for %s - %v", cacheKey, err)
		} else {
			log.Printf("API: Cached %d questions for 10 minutes", len(questions))
		}
	}

	log.Printf("API SUCCESS: Returning %d questions from database", len(questions))
	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("API ERROR: Invalid ID parameter - %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	log.Printf("API: Getting question by ID: %d", id)

	cacheKey := fmt.Sprintf("question:%d", id)

	// Пытаемся получить из кэша
	if h.cache != nil {
		var question models.Question
		err = h.cache.Get(cacheKey, &question)
		if err == nil {
			log.Printf("API: Returning question ID %d from cache", id)
			c.JSON(http.StatusOK, question)
			return
		}
		log.Printf("API: Cache miss for %s - %v", cacheKey, err)
	}

	// Если нет в кэше, то из БД
	questionPtr, err := h.storage.GetQuestionByID(id)
	if err != nil {
		log.Printf("API ERROR: Question not found ID %d - %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	question := *questionPtr

	// Оставляем в кэше на 10 минут
	if h.cache != nil {
		if err := h.cache.Set(cacheKey, question, 10*time.Minute); err != nil {
			log.Printf("API WARNING: Failed to set cache for %s - %v", cacheKey, err)
		} else {
			log.Printf("API: Cached question ID %d for 10 minutes", id)
		}
	}

	log.Printf("API SUCCESS: Returning question ID %d from database", id)
	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("API ERROR: Invalid ID parameter - %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	log.Printf("API: Deleting question ID: %d", id)

	if err := h.storage.DeleteQuestion(id); err != nil {
		log.Printf("API ERROR: Failed to delete question ID %d - %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Инвалидируем кэш
	if h.cache != nil {
		h.cache.Delete("questions:all")
		h.cache.Delete(fmt.Sprintf("question:%d", id))
		log.Printf("API: Cache invalidated for questions:all and question:%d", id)
	}

	log.Printf("API SUCCESS: Question deleted ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}
