package utils

import (
	"fmt"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type PDFExtractor struct{}

func NewPDFExtractor() *PDFExtractor {
    return &PDFExtractor{}
}

// ExtractTextFromPDF mengekstrak teks dari file PDF
func (p *PDFExtractor) ExtractTextFromPDF(filePath string) (string, error) {
    // Buka file PDF
    f, err := model.NewPdfReaderFromFile(filePath, nil)
    if err != nil {
        return "", fmt.Errorf("failed to open PDF file: %w", err)
    }

    numPages, err := f.GetNumPages()
    if err != nil {
        return "", fmt.Errorf("failed to get number of pages: %w", err)
    }

    var fullText strings.Builder

    // Ekstrak teks dari setiap halaman
    for i := 1; i <= numPages; i++ {
        page, err := f.GetPage(i)
        if err != nil {
            return "", fmt.Errorf("failed to get page %d: %w", i, err)
        }

        ex, err := extractor.New(page)
        if err != nil {
            return "", fmt.Errorf("failed to create extractor for page %d: %w", i, err)
        }

        text, err := ex.ExtractText()
        if err != nil {
            return "", fmt.Errorf("failed to extract text from page %d: %w", i, err)
        }

        fullText.WriteString(text)
        fullText.WriteString("\n\n") // Separator antar halaman
    }

    return fullText.String(), nil
}

// CleanText membersihkan teks dari karakter tidak perlu
func (p *PDFExtractor) CleanText(text string) string {
    // Hapus multiple spaces
    text = strings.Join(strings.Fields(text), " ")
    // Trim whitespace
    text = strings.TrimSpace(text)
    return text
}
