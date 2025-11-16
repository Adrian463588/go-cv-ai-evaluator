package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    ServerPort string
    OllamaURL  string
    ChromaURL  string
    UploadDir  string
}

func LoadConfig() (*Config, error) {
    // Load .env file jika ada
    godotenv.Load()

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBUser:     getEnv("DB_USER", "root"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "cv_ai_evaluator"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
        OllamaURL:  getEnv("OLLAMA_URL", "http://localhost:11434"),
        ChromaURL:  getEnv("CHROMA_URL", "http://localhost:8000"),
        UploadDir:  getEnv("UPLOAD_DIR", "./storage/uploads"),
    }

    return config, nil
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
