package handlers

import (
	"net/http"

	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"
	"cv-ai-evaluator/internal/worker"

	"github.com/gin-gonic/gin"
)

type EvaluateHandler struct {
    WorkerPool *worker.WorkerPool
}

func NewEvaluateHandler(workerPool *worker.WorkerPool) *EvaluateHandler {
    return &EvaluateHandler{WorkerPool: workerPool}
}

type EvaluateRequest struct {
    JobTitle   string `json:"job_title" binding:"required"`
    CVId       string `json:"cv_id" binding:"required"`
    ReportId   string `json:"report_id" binding:"required"`
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

    // Validasi dokumen exists
    var cvDoc, reportDoc models.UploadedDocument
    if err := database.DB.First(&cvDoc, "id = ?", req.CVId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "CV document not found"})
        return
    }

    if err := database.DB.First(&reportDoc, "id = ?", req.ReportId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Report document not found"})
        return
    }

    // Buat evaluation job
    job := models.EvaluationJob{
        CVDocumentID:      req.CVId,
        ReportDocumentID:  req.ReportId,
        JobTitleEvaluated: req.JobTitle,
        Status:            models.JobStatusQueued,
    }

    if err := database.DB.Create(&job).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create evaluation job"})
        return
    }

    // Kirim job ke worker pool
    h.WorkerPool.SubmitJob(job.ID)

    c.JSON(http.StatusOK, EvaluateResponse{
        ID:     job.ID,
        Status: string(models.JobStatusQueued),
    })
}
