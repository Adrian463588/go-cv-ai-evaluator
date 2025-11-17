package vectordb

import (
	"context"
	"fmt"
	"log"

	chromem "github.com/philippgille/chromem-go"
)

type ChromaClient struct {
	DB         *chromem.DB
	Collection *chromem.Collection
}

func NewChromaClient(persistPath string) (*ChromaClient, error) {
	// Buat embedded ChromaDB
	db := chromem.NewDB()

	// CRITICAL FIX: Gunakan model yang ADA di Ollama
	// Opsi 1: all-minilm (paling ringan, 23MB, built-in di chromem-go)
	// Opsi 2: mxbai-embed-large (harus pull dulu: ollama pull mxbai-embed-large)
	
	log.Println("⚙️  Initializing ChromaDB with embedding model...")
	
	// FIX: Gunakan all-minilm yang lebih kompatibel
	// Atau bisa pakai NewEmbeddingFuncDefault() untuk testing
	embeddingFunc := chromem.NewEmbeddingFuncOllama(
		"all-minilm",                 // Model yang lebih universal
		"http://localhost:11434",     // Base URL Ollama
	)

	log.Println("✅ Embedding function created successfully")

	// Buat atau get collection dengan Ollama embedding
	collection, err := db.GetOrCreateCollection(
		"cv_evaluator",
		map[string]string{"description": "Ground truth documents for CV evaluation"},
		embeddingFunc,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	log.Println("✅ ChromaDB collection 'cv_evaluator' ready")

	return &ChromaClient{
		DB:         db,
		Collection: collection,
	}, nil
}

// AddDocument menambahkan dokumen ke vector database
func (c *ChromaClient) AddDocument(ctx context.Context, id, content string, metadata map[string]string) error {
	err := c.Collection.AddDocument(ctx, chromem.Document{
		ID:       id,
		Content:  content,
		Metadata: metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to add document: %w", err)
	}

	return nil
}

// Query mencari dokumen yang relevan
func (c *ChromaClient) Query(ctx context.Context, queryText string, nResults int, whereFilter map[string]string) ([]chromem.Result, error) {
	results, err := c.Collection.Query(ctx, queryText, nResults, whereFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return results, nil
}

// GetRelevantContext mengambil context yang relevan untuk evaluasi
func (c *ChromaClient) GetRelevantContext(ctx context.Context, queryText string, docType string, nResults int) (string, error) {
	whereFilter := map[string]string{"type": docType}
	results, err := c.Query(ctx, queryText, nResults, whereFilter)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no relevant documents found for type: %s", docType)
	}

	// Gabungkan semua hasil menjadi satu context
	var context string
	for _, result := range results {
		context += result.Content + "\n\n"
	}

	return context, nil
}
