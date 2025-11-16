package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"cv-ai-evaluator/config"
	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
    Config *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
    return &UploadHandler{Config: cfg}
}

type UploadResponse struct {
    CVDocumentID     string `json:"cv_document_id"`
    ReportDocumentID string `json:"report_document_id"`
    Message          string `json:"message"`
}

func (h *UploadHandler) Upload(c *gin.Context) {
    // Parse multipart form
    cvFile, err := c.FormFile("cv")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "CV file is required"})
        return
    }

    reportFile, err := c.FormFile("report")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Report file is required"})
        return
    }

    // Validasi file extension
    if filepath.Ext(cvFile.Filename) != ".pdf" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "CV must be a PDF file"})
        return
    }

    if filepath.Ext(reportFile.Filename) != ".pdf" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Report must be a PDF file"})
        return
    }

    // Generate unique filenames
    cvID := uuid.New().String()
    reportID := uuid.New().String()

    cvFilename := fmt.Sprintf("%s_%s", cvID, cvFile.Filename)
    reportFilename := fmt.Sprintf("%s_%s", reportID, reportFile.Filename)

    cvPath := filepath.Join(h.Config.UploadDir, cvFilename)
    reportPath := filepath.Join(h.Config.UploadDir, reportFilename)

    // Save files
    if err := c.SaveUploadedFile(cvFile, cvPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CV file"})
        return
    }

    if err := c.SaveUploadedFile(reportFile, reportPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report file"})
        return
    }

    // Save to database
    cvDoc := models.UploadedDocument{
        ID:               cvID,
        FilePath:         cvPath,
        OriginalFilename: cvFile.Filename,
        DocumentType:     models.DocumentTypeCV,
    }

    reportDoc := models.UploadedDocument{
        ID:               reportID,
        FilePath:         reportPath,
        OriginalFilename: reportFile.Filename,
        DocumentType:     models.DocumentTypeProjectReport,
    }

    if err := database.DB.Create(&cvDoc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CV metadata"})
        return
    }

    if err := database.DB.Create(&reportDoc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report metadata"})
        return
    }

    c.JSON(http.StatusOK, UploadResponse{
        CVDocumentID:     cvID,
        ReportDocumentID: reportID,
        Message:          "Files uploaded successfully",
    })
}
