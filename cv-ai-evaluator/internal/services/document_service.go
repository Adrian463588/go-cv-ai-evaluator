package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"

	"github.com/google/uuid"
)

type DocumentService struct {
	uploadDir string
}

func NewDocumentService(uploadDir string) *DocumentService {
	return &DocumentService{
		uploadDir: uploadDir,
	}
}

// UploadDocuments handles uploading CV and Report files
func (s *DocumentService) UploadDocuments(cvFile, reportFile *multipart.FileHeader) (*models.UploadedDocument, *models.UploadedDocument, error) {
	// Validate file extensions
	if err := s.validateFileExtension(cvFile.Filename, ".pdf"); err != nil {
		return nil, nil, fmt.Errorf("CV validation failed: %w", err)
	}

	if err := s.validateFileExtension(reportFile.Filename, ".pdf"); err != nil {
		return nil, nil, fmt.Errorf("Report validation failed: %w", err)
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save CV
	cvDoc, err := s.saveFile(cvFile, models.DocumentTypeCV)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save CV: %w", err)
	}

	// Save Report
	reportDoc, err := s.saveFile(reportFile, models.DocumentTypeProjectReport)
	if err != nil {
		// Rollback CV file if report fails
		os.Remove(cvDoc.FilePath)
		database.DB.Delete(cvDoc)
		return nil, nil, fmt.Errorf("failed to save report: %w", err)
	}

	return cvDoc, reportDoc, nil
}

// saveFile saves a single file and creates database record
func (s *DocumentService) saveFile(file *multipart.FileHeader, docType models.DocumentType) (*models.UploadedDocument, error) {
	// Generate unique ID and filename
	docID := uuid.New().String()
	uniqueFilename := fmt.Sprintf("%s_%s", docID, file.Filename)
	filePath := filepath.Join(s.uploadDir, uniqueFilename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("failed to save file content: %w", err)
	}

	// Create database record
	doc := &models.UploadedDocument{
		ID:               docID,
		FilePath:         filePath,
		OriginalFilename: file.Filename,
		DocumentType:     docType,
	}

	if err := database.DB.Create(doc).Error; err != nil {
		// Rollback file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save document metadata: %w", err)
	}

	return doc, nil
}

// validateFileExtension validates file has correct extension
func (s *DocumentService) validateFileExtension(filename string, expectedExt string) error {
	ext := filepath.Ext(filename)
	if ext != expectedExt {
		return fmt.Errorf("invalid file extension: expected %s, got %s", expectedExt, ext)
	}
	return nil
}

// GetDocumentByID retrieves a document by ID
func (s *DocumentService) GetDocumentByID(id string) (*models.UploadedDocument, error) {
	var doc models.UploadedDocument
	if err := database.DB.First(&doc, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}
	return &doc, nil
}

// ValidateDocumentsExist checks if CV and Report documents exist
func (s *DocumentService) ValidateDocumentsExist(cvID, reportID string) error {
	var count int64

	// Check CV
	if err := database.DB.Model(&models.UploadedDocument{}).
		Where("id = ? AND document_type = ?", cvID, models.DocumentTypeCV).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate CV: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("CV document not found with id: %s", cvID)
	}

	// Check Report
	if err := database.DB.Model(&models.UploadedDocument{}).
		Where("id = ? AND document_type = ?", reportID, models.DocumentTypeProjectReport).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate report: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("report document not found with id: %s", reportID)
	}

	return nil
}

// DeleteDocument removes a document file and database record
func (s *DocumentService) DeleteDocument(id string) error {
	doc, err := s.GetDocumentByID(id)
	if err != nil {
		return err
	}

	// Delete file from filesystem
	if err := os.Remove(doc.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete database record
	if err := database.DB.Delete(doc).Error; err != nil {
		return fmt.Errorf("failed to delete document record: %w", err)
	}

	return nil
}

// GetDocumentsByType retrieves documents by type with pagination
func (s *DocumentService) GetDocumentsByType(docType models.DocumentType, limit, offset int) ([]models.UploadedDocument, error) {
	var docs []models.UploadedDocument
	
	query := database.DB.Where("document_type = ?", docType).
		Order("uploaded_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&docs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %w", err)
	}

	return docs, nil
}
