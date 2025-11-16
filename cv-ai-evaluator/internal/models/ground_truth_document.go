package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroundTruthType string

const (
    GroundTruthTypeJobDescription  GroundTruthType = "job_description"
    GroundTruthTypeCaseStudyBrief  GroundTruthType = "case_study_brief"
    GroundTruthTypeCVRubric        GroundTruthType = "cv_rubric"
    GroundTruthTypeProjectRubric   GroundTruthType = "project_rubric"
)

type GroundTruthDocument struct {
    ID             string          `gorm:"type:varchar(36);primaryKey" json:"id"`
    DocumentName   string          `gorm:"type:varchar(255);not null" json:"document_name"`
    DocumentType   GroundTruthType `gorm:"type:enum('job_description','case_study_brief','cv_rubric','project_rubric');not null" json:"document_type"`
    SourceFilePath string          `gorm:"type:varchar(500);not null" json:"source_file_path"`
    IngestedAt     time.Time       `gorm:"autoCreateTime" json:"ingested_at"`
    Version        string          `gorm:"type:varchar(50)" json:"version"`
}

func (GroundTruthDocument) TableName() string {
    return "ground_truth_documents"
}

func (g *GroundTruthDocument) BeforeCreate(tx *gorm.DB) error {
    if g.ID == "" {
        g.ID = uuid.New().String()
    }
    return nil
}
