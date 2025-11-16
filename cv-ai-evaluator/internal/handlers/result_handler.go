package handlers

import (
	"net/http"

	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"

	"github.com/gin-gonic/gin"
)

type ResultHandler struct{}

func NewResultHandler() *ResultHandler {
    return &ResultHandler{}
}

type ResultResponse struct {
    ID     string                 `json:"id"`
    Status string                 `json:"status"`
    Result *EvaluationResult      `json:"result,omitempty"`
    Error  string                 `json:"error,omitempty"`
}

type EvaluationResult struct {
    CVMatchRate     float64 `json:"cv_match_rate"`
    CVFeedback      string  `json:"cv_feedback"`
    ProjectScore    float64 `json:"project_score"`
    ProjectFeedback string  `json:"project_feedback"`
    OverallSummary  string  `json:"overall_summary"`
}

func (h *ResultHandler) GetResult(c *gin.Context) {
    jobID := c.Param("id")

    var job models.EvaluationJob
    if err := database.DB.First(&job, "id = ?", jobID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
        return
    }

    response := ResultResponse{
        ID:     job.ID,
        Status: string(job.Status),
    }

    switch job.Status {
    case models.JobStatusCompleted:
        response.Result = &EvaluationResult{
            CVMatchRate:     job.CVMatchRate.Float64,
            CVFeedback:      job.CVFeedback.String,
            ProjectScore:    job.ProjectScore.Float64,
            ProjectFeedback: job.ProjectFeedback.String,
            OverallSummary:  job.OverallSummary.String,
        }
    case models.JobStatusFailed:
        response.Error = job.ErrorMessage.String
    }

    c.JSON(http.StatusOK, response)
}
