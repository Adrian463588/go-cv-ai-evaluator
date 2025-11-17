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
	// Step 1: Log request details
	log.Printf("=== UPLOAD REQUEST START ===")
	log.Printf("Request Method: %s", c.Request.Method)
	log.Printf("Content-Type: %s", c.ContentType())
	log.Printf("Content-Length: %d", c.Request.ContentLength)

	// Step 2: Explicit multipart form parsing with 50MB limit
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
		log.Printf("CRITICAL: Failed to parse multipart form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid multipart form request",
			"details": err.Error(),
			"hint":    "Make sure Content-Type is multipart/form-data and files are properly attached",
		})
		return
	}

	// Step 3: Debug logging - show all available form fields
	if c.Request.MultipartForm != nil {
		if c.Request.MultipartForm.File != nil {
			keys := make([]string, 0)
			for k := range c.Request.MultipartForm.File {
				keys = append(keys, k)
			}
			log.Printf("Available file form fields: %v", keys)
		} else {
			log.Printf("WARNING: MultipartForm.File is nil - no files in request")
		}

		if c.Request.MultipartForm.Value != nil {
			valueKeys := make([]string, 0)
			for k := range c.Request.MultipartForm.Value {
				valueKeys = append(valueKeys, k)
			}
			log.Printf("Available value form fields: %v", valueKeys)
		}
	} else {
		log.Printf("CRITICAL: MultipartForm is nil")
	}

	// Step 4: Get CV file with detailed error handling
	cvFile, err := c.FormFile("cv")
	if err != nil {
		log.Printf("ERROR getting 'cv' file: %v", err)
		if errors.Is(err, http.ErrMissingFile) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Form field 'cv' is required",
				"hint":  "In Postman: Body → form-data → Key='cv', Type='File', then select your CV PDF file",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Could not process 'cv' file",
				"details": err.Error(),
			})
		}
		return
	}

	log.Printf("✅ CV file received: %s (Size: %d bytes, ContentType: %s)",
		cvFile.Filename, cvFile.Size, cvFile.Header.Get("Content-Type"))

	// Step 5: Get Report file with detailed error handling
	reportFile, err := c.FormFile("report")
	if err != nil {
		log.Printf("ERROR getting 'report' file: %v", err)
		if errors.Is(err, http.ErrMissingFile) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Form field 'report' is required",
				"hint":  "In Postman: Body → form-data → Key='report', Type='File', then select your Report PDF file",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Could not process 'report' file",
				"details": err.Error(),
			})
		}
		return
	}

	log.Printf("✅ Report file received: %s (Size: %d bytes, ContentType: %s)",
		reportFile.Filename, reportFile.Size, reportFile.Header.Get("Content-Type"))

	// Step 6: Process files
	cvDoc, reportDoc, err := h.documentService.UploadDocuments(cvFile, reportFile)
	if err != nil {
		log.Printf("Service error during upload: %v", err)
		if errors.Is(err, services.ErrInvalidFileType) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid file type provided",
				"details": err.Error(),
				"hint":    "Only PDF files are supported",
			})
			return
		}

		log.Printf("Internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process file upload",
		})
		return
	}

	// Step 7: Success response
	log.Printf("✅ Upload successful - CV: %s, Report: %s", cvDoc.ID, reportDoc.ID)
	log.Printf("=== UPLOAD REQUEST END ===")

	c.JSON(http.StatusOK, UploadResponse{
		CVDocumentID:     cvDoc.ID,
		ReportDocumentID: reportDoc.ID,
		Message:          "Files uploaded successfully",
	})
}
