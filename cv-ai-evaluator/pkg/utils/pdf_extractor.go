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
	// STEP 1: Buka file PDF
	// PERBAIKAN: Signature-nya adalah (PdfReader, io.ReadCloser, error)
	// Kita ganti nama variabel kedua menjadi 'f' (file closer)
	pdfReader, f, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	// PERBAIKAN: Tambahkan defer f.Close() untuk menghindari resource leak
	defer f.Close()

	// STEP 2: Validasi tipe isEncrypted dan check jika encrypted
	// PERBAIKAN: Panggil method .IsEncrypted() pada pdfReader
	// Method ini mengembalikan (bool, error)
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		// Tangani error jika gagal mengecek enkripsi
		return "", fmt.Errorf("failed to check encryption status: %w", err)
	}

	// PERBAIKAN: Sekarang 'isEncrypted' adalah variabel bool yang valid
	if isEncrypted {
		return "", fmt.Errorf("PDF file is encrypted and cannot be processed")
	}

	// STEP 3: Dapatkan jumlah halaman
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get number of pages: %w", err)
	}

	if numPages == 0 {
		return "", fmt.Errorf("PDF file has no pages")
	}

	// STEP 4: Ekstrak teks dari setiap halaman
	var fullText strings.Builder

	for i := 1; i <= numPages; i++ {
		// Get page
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", fmt.Errorf("failed to get page %d: %w", i, err)
		}

		// Create extractor
		ex, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("failed to create extractor for page %d: %w", i, err)
		}

		// Extract text
		text, err := ex.ExtractText()
		if err != nil {
			// Warn tapi continue, jangan fail completely
			fmt.Printf("Warning: failed to extract text from page %d: %v\n", i, err)
			continue
		}

		// Append text
		fullText.WriteString(text)
		fullText.WriteString("\n\n") // Separator antar halaman
	}

	extractedText := fullText.String()
	if extractedText == "" {
		// Ini bisa terjadi jika PDF adalah gambar atau tidak ada teks yang bisa diekstrak
		// Daripada error, kita bisa kembalikan string kosong, tapi error lebih informatif.
		return "", fmt.Errorf("no text extracted from PDF")
	}

	return extractedText, nil
}

// CleanText membersihkan teks dari karakter tidak perlu
func (p *PDFExtractor) CleanText(text string) string {
	// STEP 1: Hapus multiple spaces
	text = strings.Join(strings.Fields(text), " ")

	// STEP 2: Trim whitespace di awal dan akhir
	text = strings.TrimSpace(text)

	// STEP 3: Replace special characters
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	return text
}