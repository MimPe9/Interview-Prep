package main

import (
	"context"
	"fmt"
	"interview-prep/backend/config"
	"interview-prep/backend/internal/handlers"
	"interview-prep/backend/internal/storage"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	s := SetupStorage()
	defer s.Close()

	r := gin.Default()
	r.Use(cors.Default())

	questionHandler := handlers.NewQuestionHandler(s)

	api := r.Group("/api/v1")
	{
		api.GET("/questions", questionHandler.GetQuestions)
		api.POST("/questions", questionHandler.CreateQuestion)
	}

	r.Run(":8000")
}

func SetupStorage() *storage.PostgresStorage {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		config.DB_USER, config.DB_NAME, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_SSLMODE)

	s, err := storage.NewPosgresStorage(connStr)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	return s
}
