package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaClient struct {
    BaseURL string
    Model   string
    Client  *http.Client
}

type OllamaRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
    Options map[string]interface{} `json:"options,omitempty"`
}

type OllamaResponse struct {
    Model     string    `json:"model"`
    CreatedAt time.Time `json:"created_at"`
    Response  string    `json:"response"`
    Done      bool      `json:"done"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
    return &OllamaClient{
        BaseURL: baseURL,
        Model:   model,
        Client: &http.Client{
            Timeout: 300 * time.Second, // 5 menit timeout untuk LLM
        },
    }
}

// Generate mengirim prompt ke Ollama dan mengembalikan response
func (o *OllamaClient) Generate(prompt string, temperature float64) (string, error) {
    reqBody := OllamaRequest{
        Model:  o.Model,
        Prompt: prompt,
        Stream: false,
        Options: map[string]interface{}{
            "temperature": temperature,
        },
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    url := fmt.Sprintf("%s/api/generate", o.BaseURL)
    
    // Retry logic untuk handle timeout
    maxRetries := 3
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        resp, err := o.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", i+1, maxRetries, err)
            time.Sleep(time.Duration(i+1) * 2 * time.Second) // Exponential backoff
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bodyBytes, _ := io.ReadAll(resp.Body)
            lastErr = fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(bodyBytes))
            time.Sleep(time.Duration(i+1) * 2 * time.Second)
            continue
        }

        var ollamaResp OllamaResponse
        if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
            return "", fmt.Errorf("failed to decode response: %w", err)
        }

        return ollamaResp.Response, nil
    }

    return "", fmt.Errorf("all retry attempts failed: %w", lastErr)
}
