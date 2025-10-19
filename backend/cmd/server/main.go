package main

import (
	"fmt"
	"interview-prep/backend/config"
	"interview-prep/backend/internal/handlers"
	"interview-prep/backend/internal/storage"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	s := setupStorage()
	defer s.Close()

	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/html/index.html")

	questionHandler := handlers.NewQuestionHandler(s)

	api := r.Group("/api/v1")
	{
		api.GET("/questions", questionHandler.GetQuestions)
		api.POST("/questions", questionHandler.CreateQuestion)
	}

	r.Run(":8000")
}

func setupStorage() *storage.PostgresStorage {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		config.GetDBUser(), config.GetDBName(), config.GetDBPass(),
		config.GetDBHost(), config.GetDBPort(), config.GetDBSSLMode())

	log.Printf("Connecting to database: %s@%s:%s", config.GetDBUser(), config.GetDBHost(), config.GetDBPort())

	s, err := storage.NewPosgresStorage(connStr)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	return s
}
