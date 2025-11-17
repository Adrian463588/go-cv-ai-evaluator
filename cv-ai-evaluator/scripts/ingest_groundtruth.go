package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"cv-ai-evaluator/config"
	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"
	"cv-ai-evaluator/pkg/utils"
	"cv-ai-evaluator/pkg/vectordb"

	"github.com/google/uuid"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize ChromaDB
	chromaClient, err := vectordb.NewChromaClient("./chroma_data")
	if err != nil {
		log.Fatalf("Failed to initialize ChromaDB: %v", err)
	}

	// Initialize document reader
	docReader := utils.NewDocumentReader()
	ctx := context.Background()

	// FIX: Get absolute path for ground truth directory
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Construct proper path (works from both root and scripts directory)
	groundTruthDir := filepath.Join(workingDir, "storage", "groundtruth")
	
	// Check if we're running from scripts/ directory
	if _, err := os.Stat(groundTruthDir); os.IsNotExist(err) {
		// Try parent directory path
		groundTruthDir = filepath.Join(filepath.Dir(workingDir), "storage", "groundtruth")
	}

	log.Printf("Looking for ground truth files in: %s", groundTruthDir)

	// Verify directory exists
	if _, err := os.Stat(groundTruthDir); os.IsNotExist(err) {
		log.Fatalf("Ground truth directory not found: %s", groundTruthDir)
	}

	// Ingest documents
	documents := []struct {
		filename string
		docType  models.GroundTruthType
		docName  string
	}{
		{"job_description_backend.md", models.GroundTruthTypeJobDescription, "Backend Engineer Job Description"},
		{"case_study_brief.md", models.GroundTruthTypeCaseStudyBrief, "CV AI Evaluator Case Study"},
		{"cv_scoring_rubric.md", models.GroundTruthTypeCVRubric, "CV Evaluation Rubric"},
		{"project_scoring_rubric.md", models.GroundTruthTypeProjectRubric, "Project Evaluation Rubric"},
	}

	successCount := 0
	for _, doc := range documents {
		filePath := filepath.Join(groundTruthDir, doc.filename)

		// Check if file exists with better error message
		fileInfo, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			log.Printf("❌ File not found: %s", filePath)
			continue
		}
		if err != nil {
			log.Printf("❌ Error accessing file %s: %v", filePath, err)
			continue
		}

		log.Printf("📄 Ingesting: %s (Size: %d bytes)", doc.filename, fileInfo.Size())

		// Read document
		text, err := docReader.ReadDocument(filePath)
		if err != nil {
			log.Printf("❌ Error reading document %s: %v", doc.filename, err)
			continue
		}

		if len(text) == 0 {
			log.Printf("⚠️  Warning: %s is empty, skipping", doc.filename)
			continue
		}

		text = docReader.CleanText(text)
		log.Printf("   Extracted %d characters", len(text))

		// Generate ID
		docID := uuid.New().String()

		// Save to ChromaDB
		metadata := map[string]string{
			"type":     string(doc.docType),
			"filename": doc.filename,
			"name":     doc.docName,
		}

		if err := chromaClient.AddDocument(ctx, docID, text, metadata); err != nil {
			log.Printf("❌ Error adding document to ChromaDB: %v", err)
			continue
		}

		// Save metadata to MySQL
		gtDoc := models.GroundTruthDocument{
			ID:             docID,
			DocumentName:   doc.docName,
			DocumentType:   doc.docType,
			SourceFilePath: filePath,
			Version:        "1.0",
		}

		if err := database.DB.Create(&gtDoc).Error; err != nil {
			log.Printf("❌ Error saving to database: %v", err)
			continue
		}

		log.Printf("✅ Successfully ingested: %s (ID: %s)", doc.filename, docID)
		successCount++
	}

	log.Printf("\n🎉 Ground truth ingestion completed! (%d/%d files successful)", successCount, len(documents))
}
