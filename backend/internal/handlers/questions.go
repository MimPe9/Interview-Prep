package handlers

import (
	"fmt"
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
	var question models.Question
	if err := c.BindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.CreateQuestion(&question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Инвалидируем кэш списка вопросов
	h.cache.Delete("questions:all")

	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	cacheKey := "questions:all"

	// Пытаемся получить из кэша
	var questions []models.Question
	err := h.cache.Get(cacheKey, &questions)
	if err == nil {
		c.JSON(http.StatusOK, questions)
		return
	}

	// Если нет в кэше, то из БД
	questions, err = h.storage.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Оставляем в кэше на 10 минут
	h.cache.Set(cacheKey, questions, 10*time.Minute)

	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	cacheKey := fmt.Sprintf("question:%d", id)

	// Пытаемся получить из кэша
	var question models.Question
	err = h.cache.Get(cacheKey, &question)
	if err == nil {
		c.JSON(http.StatusOK, question)
		return
	}

	// Если нет в кэше, то из БД
	questionPtr, err := h.storage.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	question = *questionPtr

	// Оставляем в кэше на 10 минут
	h.cache.Set(cacheKey, question, 10*time.Minute)

	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.storage.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Инвалидируем кэш
	h.cache.Delete("questions:all")
	h.cache.Delete(fmt.Sprintf("question:%d", id))

	c.JSON(http.StatusOK, gin.H{"message": "Question delete"})
}
