package services

import (
	"fmt"
	"time"

	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"
)

type EvaluationService struct{}

func NewEvaluationService() *EvaluationService {
	return &EvaluationService{}
}

// CreateEvaluationJob creates a new evaluation job
func (s *EvaluationService) CreateEvaluationJob(cvID, reportID, jobTitle string) (*models.EvaluationJob, error) {
	job := &models.EvaluationJob{
		CVDocumentID:      cvID,
		ReportDocumentID:  reportID,
		JobTitleEvaluated: jobTitle,
		Status:            models.JobStatusQueued,
	}

	if err := database.DB.Create(job).Error; err != nil {
		return nil, fmt.Errorf("failed to create evaluation job: %w", err)
	}

	return job, nil
}

// GetJobByID retrieves an evaluation job by ID
func (s *EvaluationService) GetJobByID(jobID string) (*models.EvaluationJob, error) {
	var job models.EvaluationJob
	if err := database.DB.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return &job, nil
}

// GetJobWithDocuments retrieves a job with preloaded documents
func (s *EvaluationService) GetJobWithDocuments(jobID string) (*models.EvaluationJob, error) {
	var job models.EvaluationJob
	if err := database.DB.Preload("CVDocument").
		Preload("ReportDocument").
		First(&job, "id = ?", jobID).Error; err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return &job, nil
}

// UpdateJobStatus updates the status of an evaluation job
func (s *EvaluationService) UpdateJobStatus(jobID string, status models.JobStatus) error {
	return database.DB.Model(&models.EvaluationJob{}).
		Where("id = ?", jobID).
		Update("status", status).Error
}

// CompleteJob marks a job as completed with results
func (s *EvaluationService) CompleteJob(jobID string, cvMatchRate float64, cvFeedback string, projectScore float64, projectFeedback string, overallSummary string) error {
	now := time.Now()

	updates := map[string]interface{}{
		"status":           models.JobStatusCompleted,
		"completed_at":     now,
		"cv_match_rate":    cvMatchRate,
		"cv_feedback":      cvFeedback,
		"project_score":    projectScore,
		"project_feedback": projectFeedback,
		"overall_summary":  overallSummary,
	}

	if err := database.DB.Model(&models.EvaluationJob{}).
		Where("id = ?", jobID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update job results: %w", err)
	}

	return nil
}

// FailJob marks a job as failed with error message
func (s *EvaluationService) FailJob(jobID string, errorMsg string) error {
	now := time.Now()

	updates := map[string]interface{}{
		"status":        models.JobStatusFailed,
		"completed_at":  now,
		"error_message": errorMsg,
	}

	if err := database.DB.Model(&models.EvaluationJob{}).
		Where("id = ?", jobID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// GetJobsByStatus retrieves jobs by status with pagination
func (s *EvaluationService) GetJobsByStatus(status models.JobStatus, limit, offset int) ([]models.EvaluationJob, error) {
	var jobs []models.EvaluationJob

	query := database.DB.Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve jobs: %w", err)
	}

	return jobs, nil
}

// GetAllJobs retrieves all jobs with pagination
func (s *EvaluationService) GetAllJobs(limit, offset int) ([]models.EvaluationJob, int64, error) {
	var jobs []models.EvaluationJob
	var total int64

	// Count total
	if err := database.DB.Model(&models.EvaluationJob{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	// Get jobs
	query := database.DB.Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&jobs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve jobs: %w", err)
	}

	return jobs, total, nil
}

// GetJobStats returns statistics about evaluation jobs
func (s *EvaluationService) GetJobStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	statuses := []models.JobStatus{
		models.JobStatusQueued,
		models.JobStatusProcessing,
		models.JobStatusCompleted,
		models.JobStatusFailed,
	}

	for _, status := range statuses {
		var count int64
		if err := database.DB.Model(&models.EvaluationJob{}).
			Where("status = ?", status).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count status %s: %w", status, err)
		}
		stats[string(status)] = count
	}

	// Total count
	var total int64
	if err := database.DB.Model(&models.EvaluationJob{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total: %w", err)
	}
	stats["total"] = total

	return stats, nil
}

// DeleteJob removes an evaluation job
func (s *EvaluationService) DeleteJob(jobID string) error {
	if err := database.DB.Delete(&models.EvaluationJob{}, "id = ?", jobID).Error; err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}

// GetRecentJobs retrieves the most recent jobs
func (s *EvaluationService) GetRecentJobs(limit int) ([]models.EvaluationJob, error) {
	var jobs []models.EvaluationJob

	if err := database.DB.Order("created_at DESC").
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve recent jobs: %w", err)
	}

	return jobs, nil
}
