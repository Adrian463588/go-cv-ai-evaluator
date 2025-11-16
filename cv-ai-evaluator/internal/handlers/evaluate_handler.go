package handlers

import (
	"net/http"

	"cv-ai-evaluator/internal/services"
	"cv-ai-evaluator/internal/worker"

	"github.com/gin-gonic/gin"
)

type EvaluateHandler struct {
	workerPool        *worker.WorkerPool
	documentService   *services.DocumentService
	evaluationService *services.EvaluationService
}

func NewEvaluateHandler(
	workerPool *worker.WorkerPool,
	documentService *services.DocumentService,
	evaluationService *services.EvaluationService,
) *EvaluateHandler {
	return &EvaluateHandler{
		workerPool:        workerPool,
		documentService:   documentService,
		evaluationService: evaluationService,
	}
}

type EvaluateRequest struct {
	JobTitle string `json:"job_title" binding:"required"`
	CVId     string `json:"cv_id" binding:"required"`
	ReportId string `json:"report_id" binding:"required"`
}

type EvaluateResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (h *EvaluateHandler) Evaluate(c *gin.Context) {
	var req EvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate documents exist using service
	if err := h.documentService.ValidateDocumentsExist(req.CVId, req.ReportId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Create evaluation job using service
	job, err := h.evaluationService.CreateEvaluationJob(req.CVId, req.ReportId, req.JobTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Submit job to worker pool
	h.workerPool.SubmitJob(job.ID)

	c.JSON(http.StatusOK, EvaluateResponse{
		ID:     job.ID,
		Status: string(job.Status),
	})
}
