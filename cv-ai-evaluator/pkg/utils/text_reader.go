package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DocumentReader struct {
	pdfExtractor *PDFExtractor
}

func NewDocumentReader() *DocumentReader {
	return &DocumentReader{
		pdfExtractor: NewPDFExtractor(),
	}
}

// ReadDocument membaca dokumen berdasarkan extension
func (d *DocumentReader) ReadDocument(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	var rawText string
	var err error

	switch ext {
	case ".pdf":
		rawText, err = d.pdfExtractor.ExtractTextFromPDF(filePath)
	case ".md", ".txt":
		rawText, err = d.readTextFile(filePath)
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return "", err
	}

	// BUG FIX: Panggil CleanText SEBELUM mengembalikan
	return d.CleanText(rawText), nil
}

// readTextFile membaca file text/markdown
func (d *DocumentReader) readTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// CleanText membersihkan teks
func (d *DocumentReader) CleanText(text string) string {
	// STEP 1: Hapus multiple spaces
	text = strings.Join(strings.Fields(text), " ")
	
	// STEP 2: Trim whitespace di awal dan akhir
	text = strings.TrimSpace(text)
	
	// STEP 3: (FIX) Standarisasi newlines (agar konsisten dgn pdf_extractor)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	
	return text
}