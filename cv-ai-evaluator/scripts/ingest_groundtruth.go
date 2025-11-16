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

    // Initialize document reader (support PDF, MD, TXT)
    docReader := utils.NewDocumentReader()

    ctx := context.Background()

    // Directory untuk ground truth documents
    groundTruthDir := "./storage/groundtruth"

    // Ingest documents - BISA PDF atau MD atau TXT
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

    for _, doc := range documents {
        filePath := filepath.Join(groundTruthDir, doc.filename)

        // Check if file exists
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            log.Printf("Warning: File not found: %s, skipping...", filePath)
            continue
        }

        log.Printf("Ingesting: %s", doc.filename)

        // Read document (support PDF, MD, TXT)
        text, err := docReader.ReadDocument(filePath)
        if err != nil {
            log.Printf("Error reading document %s: %v", doc.filename, err)
            continue
        }

        text = docReader.CleanText(text)

        // Generate ID
        docID := uuid.New().String()

        // Save to ChromaDB
        metadata := map[string]string{
            "type":     string(doc.docType),
            "filename": doc.filename,
            "name":     doc.docName,
        }

        if err := chromaClient.AddDocument(ctx, docID, text, metadata); err != nil {
            log.Printf("Error adding document to ChromaDB: %v", err)
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
            log.Printf("Error saving to database: %v", err)
            continue
        }

        log.Printf("✓ Successfully ingested: %s", doc.filename)
    }

    log.Println("Ground truth ingestion completed!")
}
