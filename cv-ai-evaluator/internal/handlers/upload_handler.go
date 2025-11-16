package handlers

import (
	"net/http"

	"cv-ai-evaluator/internal/services"

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

	// Use service to handle upload
	cvDoc, reportDoc, err := h.documentService.UploadDocuments(cvFile, reportFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		CVDocumentID:     cvDoc.ID,
		ReportDocumentID: reportDoc.ID,
		Message:          "Files uploaded successfully",
	})
}
