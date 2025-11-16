package services

import (
	"errors" // PERBAIKAN: Tambahkan import errors
	"fmt"
	"io" // PERBAIKAN: Tambahkan import io
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"

	"github.com/google/uuid"
)

// PERBAIKAN: Definisi Sentinel Errors untuk error handling yang bersih
var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileReadError   = errors.New("could not read uploaded file")
	ErrFileSaveError   = errors.New("could not save file to disk")
	ErrDatabaseError   = errors.New("database operation failed")
)

type DocumentService struct {
	uploadDir string
}

func NewDocumentService(uploadDir string) *DocumentService {
	return &DocumentService{
		uploadDir: uploadDir,
	}
}

// UploadDocuments menangani upload file CV dan Report
func (s *DocumentService) UploadDocuments(cvFile, reportFile *multipart.FileHeader) (*models.UploadedDocument, *models.UploadedDocument, error) {
	// Validasi tipe file
	if err := s.validateFileExtension(cvFile.Filename, ".pdf"); err != nil {
		return nil, nil, fmt.Errorf("CV validation failed: %w", err)
	}

	if err := s.validateFileExtension(reportFile.Filename, ".pdf"); err != nil {
		return nil, nil, fmt.Errorf("report validation failed: %w", err)
	}

	// Buat direktori upload jika belum ada
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Simpan CV
	cvDoc, err := s.saveFile(cvFile, models.DocumentTypeCV)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save CV: %w", err)
	}

	// Simpan Report
	reportDoc, err := s.saveFile(reportFile, models.DocumentTypeProjectReport)
	if err != nil {
		// Rollback (hapus file CV) jika file report gagal disimpan
		if removeErr := os.Remove(cvDoc.FilePath); removeErr != nil {
			log.Printf("Warning: failed to rollback file %s: %v", cvDoc.FilePath, removeErr)
		}
		if dbErr := database.DB.Delete(cvDoc).Error; dbErr != nil {
			log.Printf("Warning: failed to rollback db record %s: %v", cvDoc.ID, dbErr)
		}
		return nil, nil, fmt.Errorf("failed to save report: %w", err)
	}

	return cvDoc, reportDoc, nil
}

// saveFile menyimpan satu file dan membuat record di database
func (s *DocumentService) saveFile(file *multipart.FileHeader, docType models.DocumentType) (*models.UploadedDocument, error) {
	// Buat ID unik dan nama file
	docID := uuid.New().String()
	// PERBAIKAN: Bersihkan nama file menggunakan filepath.Base untuk keamanan
	uniqueFilename := fmt.Sprintf("%s_%s", docID, filepath.Base(file.Filename))
	filePath := filepath.Join(s.uploadDir, uniqueFilename)

	// Buka file yang di-upload
	src, err := file.Open()
	if err != nil {
		// PERBAIKAN: Bungkus (wrap) error dengan sentinel error
		return nil, fmt.Errorf("%w: %v", ErrFileReadError, err)
	}
	defer src.Close()

	// Buat file tujuan
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileSaveError, err)
	}
	defer dst.Close()

	// Copy konten file
	if _, err := io.Copy(dst, src); err != nil { // PERBAIKAN: Gunakan io.Copy
		return nil, fmt.Errorf("%w: %v", ErrFileSaveError, err)
	}

	// Buat record di database
	doc := &models.UploadedDocument{
		ID:               docID,
		FilePath:         filePath,
		OriginalFilename: file.Filename,
		DocumentType:     docType,
	}

	if err := database.DB.Create(doc).Error; err != nil {
		// Rollback (hapus file) jika insert ke DB gagal
		if removeErr := os.Remove(filePath); removeErr != nil {
			log.Printf("Warning: failed to rollback file %s on db error: %v", filePath, removeErr)
		}
		// PERBAIKAN: Bungkus error database
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return doc, nil
}

// validateFileExtension memvalidasi file memiliki ekstensi yang benar
func (s *DocumentService) validateFileExtension(filename string, expectedExt string) error {
	ext := filepath.Ext(filename)
	if ext != expectedExt {
		// PERBAIKAN: Kembalikan error yang dibungkus dengan sentinel error
		return fmt.Errorf("%w: expected %s, got %s", ErrInvalidFileType, expectedExt, ext)
	}
	return nil
}

// GetDocumentByID mengambil dokumen berdasarkan ID
func (s *DocumentService) GetDocumentByID(id string) (*models.UploadedDocument, error) {
	var doc models.UploadedDocument
	if err := database.DB.First(&doc, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}
	return &doc, nil
}

// ValidateDocumentsExist memeriksa apakah dokumen CV dan Report ada
func (s *DocumentService) ValidateDocumentsExist(cvID, reportID string) error {
	var count int64

	// Cek CV
	if err := database.DB.Model(&models.UploadedDocument{}).
		Where("id = ? AND document_type = ?", cvID, models.DocumentTypeCV).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate CV: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("CV document not found with id: %s", cvID)
	}

	// Cek Report
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