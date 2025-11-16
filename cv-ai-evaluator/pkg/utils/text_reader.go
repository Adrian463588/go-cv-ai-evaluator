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
    
    switch ext {
    case ".pdf":
        return d.pdfExtractor.ExtractTextFromPDF(filePath)
    case ".md", ".txt":
        return d.readTextFile(filePath)
    default:
        return "", fmt.Errorf("unsupported file format: %s", ext)
    }
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
    // Hapus multiple spaces
    text = strings.Join(strings.Fields(text), " ")
    // Trim whitespace
    text = strings.TrimSpace(text)
    return text
}
