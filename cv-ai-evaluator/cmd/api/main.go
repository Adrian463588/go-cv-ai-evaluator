package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"cv-ai-evaluator/config"
	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/handlers"
	"cv-ai-evaluator/internal/services"
	"cv-ai-evaluator/internal/worker"
	"cv-ai-evaluator/pkg/llm"
	"cv-ai-evaluator/pkg/vectordb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize Ollama client
	ollamaClient := llm.NewOllamaClient(cfg.OllamaURL, "gemma3:4b")

	// Initialize ChromaDB client
	chromaClient, err := vectordb.NewChromaClient("./chroma_data")
	if err != nil {
		log.Fatalf("Failed to initialize ChromaDB: %v", err)
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Initialize services
	documentService := services.NewDocumentService(cfg.UploadDir)
	evaluationService := services.NewEvaluationService()

	// Initialize worker pool with services
	workerPool := worker.NewWorkerPool(3, ollamaClient, chromaClient, evaluationService)
	workerPool.Start()

	// Setup Gin router
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handlers with services
	uploadHandler := handlers.NewUploadHandler(documentService)
	evaluateHandler := handlers.NewEvaluateHandler(workerPool, documentService, evaluationService)
	resultHandler := handlers.NewResultHandler(evaluationService)

	// Routes
	router.POST("/upload", uploadHandler.Upload)
	router.POST("/evaluate", evaluateHandler.Evaluate)
	router.GET("/result/:id", resultHandler.GetResult)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gracefully...")
		workerPool.Stop()
		os.Exit(0)
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
