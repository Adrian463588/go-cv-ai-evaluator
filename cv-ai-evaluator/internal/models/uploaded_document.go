package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentType string

const (
    DocumentTypeCV            DocumentType = "cv"
    DocumentTypeProjectReport DocumentType = "project_report"
)

type UploadedDocument struct {
    ID               string       `gorm:"type:varchar(36);primaryKey" json:"id"`
    FilePath         string       `gorm:"type:varchar(500);not null" json:"file_path"`
    OriginalFilename string       `gorm:"type:varchar(255);not null" json:"original_filename"`
    DocumentType     DocumentType `gorm:"type:enum('cv','project_report');not null" json:"document_type"`
    UploadedAt       time.Time    `gorm:"autoCreateTime" json:"uploaded_at"`
}

func (UploadedDocument) TableName() string {
    return "uploaded_documents"
}

// Hook sebelum create untuk generate UUID
func (u *UploadedDocument) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    return nil
}
