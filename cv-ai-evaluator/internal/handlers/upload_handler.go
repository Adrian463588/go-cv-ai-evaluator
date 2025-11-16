package handlers

import (
	"cv-ai-evaluator/internal/services"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
    documentService *services.DocumentService
}

func NewUploadHandler(documentService *services.DocumentService) *UploadHandler {
    return &UploadHandler{
        documentService: documentService,
    }
}

type UploadResponse struct {
    CVDocumentID     string `json:"cv_document_id"`
    ReportDocumentID string `json:"report_document_id"`
    Message          string `json:"message"`
}

func (h *UploadHandler) Upload(c *gin.Context) {
    // Step 1: Explicit multipart form parsing with 50MB limit
    if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
        log.Printf("CRITICAL: Failed to parse multipart form: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid multipart form request",
        })
        return
    }

    // Step 2: Debug logging
    log.Printf("Incoming request - ContentType: %s", c.ContentType())
    if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
        keys := make([]string, 0)
        for k := range c.Request.MultipartForm.File {
            keys = append(keys, k)
        }
        log.Printf("Available form fields: %v", keys)
    }

    // Step 3: Get CV file
    cvFile, err := c.FormFile("cv")
    if err != nil {
        log.Printf("ERROR getting 'cv' file: %v", err)
        if errors.Is(err, http.ErrMissingFile) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Form field 'cv' is required"})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Could not process 'cv' file"})
        }
        return
    }
    log.Printf("CV file received: %s (Size: %d bytes)", cvFile.Filename, cvFile.Size)

    // Step 4: Get Report file
    reportFile, err := c.FormFile("report")
    if err != nil {
        log.Printf("ERROR getting 'report' file: %v", err)
        if errors.Is(err, http.ErrMissingFile) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Form field 'report' is required"})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Could not process 'report' file"})
        }
        return
    }
    log.Printf("Report file received: %s (Size: %d bytes)", reportFile.Filename, reportFile.Size)

    // Step 5: Process files
    cvDoc, reportDoc, err := h.documentService.UploadDocuments(cvFile, reportFile)
    if err != nil {
        log.Printf("Service error during upload: %v", err)
        if errors.Is(err, services.ErrInvalidFileType) {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid file type provided",
                "details": err.Error(),
            })
            return
        }

        log.Printf("Internal server error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to process file upload",
        })
        return
    }

    // Step 6: Success response
    c.JSON(http.StatusOK, UploadResponse{
        CVDocumentID:     cvDoc.ID,
        ReportDocumentID: reportDoc.ID,
        Message:          "Files uploaded successfully",
    })
}
