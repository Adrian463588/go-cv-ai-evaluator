package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobStatus string

const (
    JobStatusQueued     JobStatus = "queued"
    JobStatusProcessing JobStatus = "processing"
    JobStatusCompleted  JobStatus = "completed"
    JobStatusFailed     JobStatus = "failed"
)

type EvaluationJob struct {
    ID                 string         `gorm:"type:varchar(36);primaryKey" json:"id"`
    CVDocumentID       string         `gorm:"type:varchar(36);not null" json:"cv_document_id"`
    ReportDocumentID   string         `gorm:"type:varchar(36);not null" json:"report_document_id"`
    JobTitleEvaluated  string         `gorm:"type:varchar(255);not null" json:"job_title_evaluated"`
    Status             JobStatus      `gorm:"type:enum('queued','processing','completed','failed');default:'queued'" json:"status"`
    CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
    CompletedAt        sql.NullTime   `json:"completed_at,omitempty"`
    ErrorMessage       sql.NullString `gorm:"type:text" json:"error_message,omitempty"`
    CVMatchRate        sql.NullFloat64 `gorm:"type:decimal(3,2)" json:"cv_match_rate,omitempty"`
    CVFeedback         sql.NullString `gorm:"type:text" json:"cv_feedback,omitempty"`
    ProjectScore       sql.NullFloat64 `gorm:"type:decimal(3,2)" json:"project_score,omitempty"`
    ProjectFeedback    sql.NullString `gorm:"type:text" json:"project_feedback,omitempty"`
    OverallSummary     sql.NullString `gorm:"type:text" json:"overall_summary,omitempty"`

    // Relations
    CVDocument     UploadedDocument `gorm:"foreignKey:CVDocumentID" json:"-"`
    ReportDocument UploadedDocument `gorm:"foreignKey:ReportDocumentID" json:"-"`
}

func (EvaluationJob) TableName() string {
    return "evaluation_jobs"
}

func (e *EvaluationJob) BeforeCreate(tx *gorm.DB) error {
    if e.ID == "" {
        e.ID = uuid.New().String()
    }
    return nil
}
