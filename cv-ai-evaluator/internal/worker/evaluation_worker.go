package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"cv-ai-evaluator/internal/database"
	"cv-ai-evaluator/internal/models"
	"cv-ai-evaluator/pkg/llm"
	"cv-ai-evaluator/pkg/utils"
	"cv-ai-evaluator/pkg/vectordb"
)

type WorkerPool struct {
    jobQueue      chan string
    workerCount   int
    wg            sync.WaitGroup
    ctx           context.Context
    cancel        context.CancelFunc
    ollamaClient  *llm.OllamaClient
    chromaClient  *vectordb.ChromaClient
    pdfExtractor  *utils.PDFExtractor
}

func NewWorkerPool(workerCount int, ollamaClient *llm.OllamaClient, chromaClient *vectordb.ChromaClient) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &WorkerPool{
        jobQueue:     make(chan string, 100), // Buffer untuk 100 jobs
        workerCount:  workerCount,
        ctx:          ctx,
        cancel:       cancel,
        ollamaClient: ollamaClient,
        chromaClient: chromaClient,
        pdfExtractor: utils.NewPDFExtractor(),
    }
}

// Start memulai worker pool
func (wp *WorkerPool) Start() {
    for i := 1; i <= wp.workerCount; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
    log.Printf("Started %d workers", wp.workerCount)
}

// Stop menghentikan worker pool
func (wp *WorkerPool) Stop() {
    log.Println("Stopping worker pool...")
    close(wp.jobQueue)
    wp.cancel()
    wp.wg.Wait()
    log.Println("Worker pool stopped")
}

// SubmitJob menambahkan job ke queue
func (wp *WorkerPool) SubmitJob(jobID string) {
    wp.jobQueue <- jobID
}

// worker adalah goroutine yang memproses jobs
func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    log.Printf("Worker %d started", id)
    
    for {
        select {
        case <-wp.ctx.Done():
            log.Printf("Worker %d stopping due to context cancellation", id)
            return
        case jobID, ok := <-wp.jobQueue:
            if !ok {
                log.Printf("Worker %d stopping: job queue closed", id)
                return
            }
            
            log.Printf("Worker %d processing job: %s", id, jobID)
            if err := wp.processJob(jobID); err != nil {
                log.Printf("Worker %d failed to process job %s: %v", id, jobID, err)
            } else {
                log.Printf("Worker %d completed job: %s", id, jobID)
            }
        }
    }
}

// processJob memproses satu evaluation job
func (wp *WorkerPool) processJob(jobID string) error {
    // 1. Update status ke processing
    if err := database.DB.Model(&models.EvaluationJob{}).
        Where("id = ?", jobID).
        Update("status", models.JobStatusProcessing).Error; err != nil {
        return fmt.Errorf("failed to update job status: %w", err)
    }

    // 2. Load job data
    var job models.EvaluationJob
    if err := database.DB.Preload("CVDocument").Preload("ReportDocument").
        First(&job, "id = ?", jobID).Error; err != nil {
        return wp.markJobAsFailed(jobID, fmt.Sprintf("failed to load job: %v", err))
    }

    // 3. Extract PDF text
    cvText, err := wp.pdfExtractor.ExtractTextFromPDF(job.CVDocument.FilePath)
    if err != nil {
        return wp.markJobAsFailed(jobID, fmt.Sprintf("failed to extract CV text: %v", err))
    }
    cvText = wp.pdfExtractor.CleanText(cvText)

    reportText, err := wp.pdfExtractor.ExtractTextFromPDF(job.ReportDocument.FilePath)
    if err != nil {
        return wp.markJobAsFailed(jobID, fmt.Sprintf("failed to extract report text: %v", err))
    }
    reportText = wp.pdfExtractor.CleanText(reportText)

    // 4. Evaluate CV
    cvMatchRate, cvFeedback, err := wp.evaluateCV(cvText, job.JobTitleEvaluated)
    if err != nil {
        return wp.markJobAsFailed(jobID, fmt.Sprintf("CV evaluation failed: %v", err))
    }

    // 5. Evaluate Project Report
    projectScore, projectFeedback, err := wp.evaluateProject(reportText)
    if err != nil {
        return wp.markJobAsFailed(jobID, fmt.Sprintf("Project evaluation failed: %v", err))
    }

    // 6. Generate Overall Summary
    overallSummary, err := wp.generateOverallSummary(cvMatchRate, cvFeedback, projectScore, projectFeedback, job.JobTitleEvaluated)
    if err != nil {
        return wp.markJobAsFailed(jobID, fmt.Sprintf("Summary generation failed: %v", err))
    }

    // 7. Update job dengan hasil
    now := time.Now()
    updates := map[string]interface{}{
        "status":           models.JobStatusCompleted,
        "completed_at":     now,
        "cv_match_rate":    cvMatchRate,
        "cv_feedback":      cvFeedback,
        "project_score":    projectScore,
        "project_feedback": projectFeedback,
        "overall_summary":  overallSummary,
    }

    if err := database.DB.Model(&models.EvaluationJob{}).
        Where("id = ?", jobID).
        Updates(updates).Error; err != nil {
        return fmt.Errorf("failed to save results: %w", err)
    }

    return nil
}

// evaluateCV melakukan evaluasi CV dengan RAG dan LLM
func (wp *WorkerPool) evaluateCV(cvText, jobTitle string) (float64, string, error) {
    // 1. Query vector DB untuk job description dan rubric
    jdContext, err := wp.chromaClient.GetRelevantContext(
        wp.ctx,
        fmt.Sprintf("%s job description requirements", jobTitle),
        "job_description",
        2,
    )
    if err != nil {
        log.Printf("Warning: failed to get job description context: %v", err)
        jdContext = "No specific job description available."
    }

    rubricContext, err := wp.chromaClient.GetRelevantContext(
        wp.ctx,
        "CV evaluation rubric scoring criteria",
        "cv_rubric",
        1,
    )
    if err != nil {
        log.Printf("Warning: failed to get CV rubric context: %v", err)
        rubricContext = "Evaluate based on standard criteria."
    }

    // 2. Build prompt untuk LLM
    prompt := fmt.Sprintf(`You are an expert technical recruiter evaluating a candidate's CV for a %s position.

Job Description and Requirements:
%s

Evaluation Rubric:
%s

Candidate's CV:
%s

Based on the job requirements and evaluation rubric, please:
1. Rate the overall CV match (0.0 to 1.0 scale) based on:
   - Technical skills match (40%%)
   - Experience level (25%%)
   - Relevant achievements (20%%)
   - Cultural fit indicators (15%%)
   
2. Provide detailed feedback (3-5 sentences) covering strengths, gaps, and recommendations.

IMPORTANT: Your response MUST be valid JSON in this exact format:
{
  "match_rate": 0.82,
  "feedback": "Your detailed feedback here..."
}`, jobTitle, jdContext, rubricContext, cvText)

    // 3. Call LLM
    response, err := wp.ollamaClient.Generate(prompt, 0.3) // Low temperature untuk consistency
    if err != nil {
        return 0, "", fmt.Errorf("LLM call failed: %w", err)
    }

    // 4. Parse response
    matchRate, feedback, err := wp.parseCVResponse(response)
    if err != nil {
        return 0, "", fmt.Errorf("failed to parse CV response: %w", err)
    }

    return matchRate, feedback, nil
}

// evaluateProject melakukan evaluasi project report
func (wp *WorkerPool) evaluateProject(reportText string) (float64, string, error) {
    // 1. Query vector DB untuk case study brief dan rubric
    briefContext, err := wp.chromaClient.GetRelevantContext(
        wp.ctx,
        "case study brief requirements specifications",
        "case_study_brief",
        1,
    )
    if err != nil {
        log.Printf("Warning: failed to get case study brief: %v", err)
        briefContext = "Evaluate based on general backend project standards."
    }

    rubricContext, err := wp.chromaClient.GetRelevantContext(
        wp.ctx,
        "project evaluation rubric scoring criteria",
        "project_rubric",
        1,
    )
    if err != nil {
        log.Printf("Warning: failed to get project rubric: %v", err)
        rubricContext = "Evaluate based on standard project criteria."
    }

    // 2. Build prompt
    prompt := fmt.Sprintf(`You are an expert technical evaluator reviewing a candidate's project report.

Case Study Requirements:
%s

Evaluation Rubric:
%s

Candidate's Project Report:
%s

Based on the requirements and rubric, please:
1. Provide a score (1.0 to 5.0 scale) based on:
   - Correctness & completeness (30%%)
   - Code quality & structure (25%%)
   - Resilience & error handling (20%%)
   - Documentation quality (15%%)
   - Creativity & extras (10%%)
   
2. Provide detailed feedback (3-5 sentences) on strengths, weaknesses, and improvements.

IMPORTANT: Your response MUST be valid JSON in this exact format:
{
  "score": 4.2,
  "feedback": "Your detailed feedback here..."
}`, briefContext, rubricContext, reportText)

    // 3. Call LLM
    response, err := wp.ollamaClient.Generate(prompt, 0.3)
    if err != nil {
        return 0, "", fmt.Errorf("LLM call failed: %w", err)
    }

    // 4. Parse response
    score, feedback, err := wp.parseProjectResponse(response)
    if err != nil {
        return 0, "", fmt.Errorf("failed to parse project response: %w", err)
    }

    return score, feedback, nil
}

// generateOverallSummary membuat ringkasan keseluruhan
func (wp *WorkerPool) generateOverallSummary(cvMatchRate float64, cvFeedback string, projectScore float64, projectFeedback, jobTitle string) (string, error) {
    prompt := fmt.Sprintf(`You are a senior technical hiring manager making a final decision on a candidate for a %s position.

CV Evaluation:
- Match Rate: %.2f (0-1 scale)
- Feedback: %s

Project Evaluation:
- Score: %.1f (1-5 scale)
- Feedback: %s

Based on both evaluations, provide a 3-5 sentence overall summary that:
1. Summarizes the candidate's strengths
2. Identifies key gaps or concerns
3. Gives a hiring recommendation (strong hire / hire / maybe / no hire)

Be direct, professional, and actionable.`, jobTitle, cvMatchRate, cvFeedback, projectScore, projectFeedback)

    response, err := wp.ollamaClient.Generate(prompt, 0.4)
    if err != nil {
        return "", fmt.Errorf("LLM call failed: %w", err)
    }

    return strings.TrimSpace(response), nil
}

// parseCVResponse mem-parse JSON response dari LLM untuk CV evaluation
func (wp *WorkerPool) parseCVResponse(response string) (float64, string, error) {
    // Coba parse sebagai JSON
    var result struct {
        MatchRate float64 `json:"match_rate"`
        Feedback  string  `json:"feedback"`
    }

    // Cari JSON block dalam response
    jsonStart := strings.Index(response, "{")
    jsonEnd := strings.LastIndex(response, "}")
    
    if jsonStart >= 0 && jsonEnd > jsonStart {
        jsonStr := response[jsonStart : jsonEnd+1]
        if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
            return result.MatchRate, result.Feedback, nil
        }
    }

    // Fallback: extract dengan regex
    matchRateRegex := regexp.MustCompile(`"?match_rate"?\s*:?\s*([0-9.]+)`)
    feedbackRegex := regexp.MustCompile(`"?feedback"?\s*:?\s*"([^"]+)"`)

    matchRateMatch := matchRateRegex.FindStringSubmatch(response)
    feedbackMatch := feedbackRegex.FindStringSubmatch(response)

    if len(matchRateMatch) > 1 && len(feedbackMatch) > 1 {
        matchRate, err := strconv.ParseFloat(matchRateMatch[1], 64)
        if err != nil {
            return 0, "", fmt.Errorf("failed to parse match_rate: %w", err)
        }
        return matchRate, feedbackMatch[1], nil
    }

    return 0, "", fmt.Errorf("could not parse CV response: %s", response)
}

// parseProjectResponse mem-parse JSON response dari LLM untuk project evaluation
func (wp *WorkerPool) parseProjectResponse(response string) (float64, string, error) {
    var result struct {
        Score    float64 `json:"score"`
        Feedback string  `json:"feedback"`
    }

    // Cari JSON block
    jsonStart := strings.Index(response, "{")
    jsonEnd := strings.LastIndex(response, "}")
    
    if jsonStart >= 0 && jsonEnd > jsonStart {
        jsonStr := response[jsonStart : jsonEnd+1]
        if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
            return result.Score, result.Feedback, nil
        }
    }

    // Fallback: extract dengan regex
    scoreRegex := regexp.MustCompile(`"?score"?\s*:?\s*([0-9.]+)`)
    feedbackRegex := regexp.MustCompile(`"?feedback"?\s*:?\s*"([^"]+)"`)

    scoreMatch := scoreRegex.FindStringSubmatch(response)
    feedbackMatch := feedbackRegex.FindStringSubmatch(response)

    if len(scoreMatch) > 1 && len(feedbackMatch) > 1 {
        score, err := strconv.ParseFloat(scoreMatch[1], 64)
        if err != nil {
            return 0, "", fmt.Errorf("failed to parse score: %w", err)
        }
        return score, feedbackMatch[1], nil
    }

    return 0, "", fmt.Errorf("could not parse project response: %s", response)
}

// markJobAsFailed menandai job sebagai failed
func (wp *WorkerPool) markJobAsFailed(jobID, errorMsg string) error {
    now := time.Now()
    return database.DB.Model(&models.EvaluationJob{}).
        Where("id = ?", jobID).
        Updates(map[string]interface{}{
            "status":        models.JobStatusFailed,
            "completed_at":  now,
            "error_message": sql.NullString{String: errorMsg, Valid: true},
        }).Error
}
