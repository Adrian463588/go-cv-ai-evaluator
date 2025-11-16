package handlers

import (
	"net/http"

	"cv-ai-evaluator/internal/models"
	"cv-ai-evaluator/internal/services"

	"github.com/gin-gonic/gin"
)

type ResultHandler struct {
	evaluationService *services.EvaluationService
}

func NewResultHandler(evaluationService *services.EvaluationService) *ResultHandler {
	return &ResultHandler{
		evaluationService: evaluationService,
	}
}

type ResultResponse struct {
	ID     string            `json:"id"`
	Status string            `json:"status"`
	Result *EvaluationResult `json:"result,omitempty"`
	Error  string            `json:"error,omitempty"`
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

	// Get job using service
	job, err := h.evaluationService.GetJobByID(jobID)
	if err != nil {
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
