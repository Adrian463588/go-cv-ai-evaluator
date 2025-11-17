# Project Scoring Rubric

## 1. CV Match Evaluation (1–5 scale per parameter)

### Technical Skills Match (Weight: 40%)

**Description:** Alignment with job requirements (backend, databases, APIs, cloud, AI/LLM)

**Scoring Guide:**
- **1** = Irrelevant skills
- **2** = Few overlaps
- **3** = Partial match
- **4** = Strong match
- **5** = Excellent match + AI/LLM exposure

### Experience Level (Weight: 25%)

**Description:** Years of experience and project complexity

**Scoring Guide:**
- **1** = <1 yr / trivial projects
- **2** = 1–2 yrs
- **3** = 2–3 yrs with mid-scale projects
- **4** = 3–4 yrs solid track record
- **5** = 5+ yrs / high-impact projects

### Relevant Achievements (Weight: 20%)

**Description:** Impact of past work (scaling, performance, adoption)

**Scoring Guide:**
- **1** = No clear achievements
- **2** = Minimal improvements
- **3** = Some measurable outcomes
- **4** = Significant contributions
- **5** = Major measurable impact

### Cultural / Collaboration Fit (Weight: 15%)

**Description:** Communication, learning mindset, teamwork/leadership

**Scoring Guide:**
- **1** = Not demonstrated
- **2** = Minimal
- **3** = Average
- **4** = Good
- **5** = Excellent and well-demonstrated

---

## 2. Project Deliverable Evaluation (1–5 scale per parameter)

### Correctness (Prompt & Chaining) (Weight: 30%)

**Description:** Implements prompt design, LLM chaining, RAG context injection

**Scoring Guide:**
- **1** = Not implemented
- **2** = Minimal attempt
- **3** = Works partially
- **4** = Works correctly
- **5** = Fully correct + thoughtful

### Code Quality & Structure (Weight: 25%)

**Description:** Clean, modular, reusable, tested

**Scoring Guide:**
- **1** = Poor
- **2** = Some structure
- **3** = Decent modularity
- **4** = Good structure + some tests
- **5** = Excellent quality + strong tests

### Resilience & Error Handling (Weight: 20%)

**Description:** Handles long jobs, retries, randomness, API failures

**Scoring Guide:**
- **1** = Missing
- **2** = Minimal
- **3** = Partial handling
- **4** = Solid handling
- **5** = Robust, production-ready

### Documentation & Explanation (Weight: 15%)

**Description:** README clarity, setup instructions, trade-off explanations

**Scoring Guide:**
- **1** = Missing
- **2** = Minimal
- **3** = Adequate
- **4** = Clear
- **5** = Excellent + insightful

### Creativity / Bonus (Weight: 10%)

**Description:** Extra features beyond requirements

**Scoring Guide:**
- **1** = None
- **2** = Very basic
- **3** = Useful extras
- **4** = Strong enhancements
- **5** = Outstanding creativity

---

## 3. Overall Candidate Evaluation

### CV Match Rate Calculation

**Formula:** Weighted Average (1–5) → Convert to 0-1 decimal (×0.2)

**Example Calculation:**
```
Technical Skills Match: 4 × 0.40 = 1.60
Experience Level: 3 × 0.25 = 0.75
Relevant Achievements: 4 × 0.20 = 0.80
Cultural Fit: 4 × 0.15 = 0.60
-----------------------------------
Total Weighted Score: 3.75
CV Match Rate: 3.75 × 0.20 = 0.75
```

### Project Score Calculation

**Formula:** Weighted Average (1–5)

**Example Calculation:**
```
Correctness: 4 × 0.30 = 1.20
Code Quality: 4 × 0.25 = 1.00
Resilience: 3 × 0.20 = 0.60
Documentation: 4 × 0.15 = 0.60
Creativity: 3 × 0.10 = 0.30
-----------------------------------
Project Score: 3.70
```

### Overall Summary Requirements

**Service should return 3–5 sentences covering:**
- **Strengths:** Key positive aspects from both CV and project evaluation
- **Gaps:** Areas for improvement or concerns identified
- **Recommendations:** Hiring decision or next steps (e.g., "Strong hire", "Hire with conditions", "Further interview needed", "Not recommended")

**Example:**
> "The candidate demonstrates strong backend fundamentals with solid experience in Node.js and database design, showing a CV match rate of 0.75. The project implementation (score: 3.7) correctly implements prompt chaining and RAG retrieval, though error handling could be more robust. Overall, this is a solid hire for a mid-level backend role, with potential for growth in AI/LLM integration. Recommend moving forward to technical interview."

---

## Evaluation Process Flow

```
1. Receive CV and Project Report
2. Extract and parse content from PDFs
3. Query vector database for relevant context:
   - Job descriptions for CV evaluation
   - Case study brief for project evaluation
   - Respective rubrics for both
4. LLM Chain Evaluation:
   a. CV Evaluation → cv_match_rate, cv_feedback
   b. Project Evaluation → project_score, project_feedback
   c. Final Synthesis → overall_summary
5. Store and return results
```

---

## Quality Assurance Guidelines

### For Evaluators

- **Consistency:** Use temperature control (0.3-0.4) for LLM calls to ensure stable scoring
- **Context Retrieval:** Verify that relevant sections from ground truth documents are properly retrieved
- **Score Validation:** Ensure scores fall within expected ranges (0-1 for match_rate, 1-5 for project_score)
- **Feedback Quality:** Check that feedback is specific, actionable, and references actual content from documents

### For Implementation

- **Retry Logic:** Implement exponential backoff for LLM API failures
- **Timeout Handling:** Set appropriate timeouts (5+ minutes for complex evaluations)
- **Error Logging:** Capture and log all errors for debugging
- **Result Validation:** Verify JSON structure and data types before storing results

---

## Score Interpretation

### CV Match Rate (0.0 - 1.0)

- **0.00 - 0.40:** Poor fit - Missing critical skills or experience
- **0.41 - 0.60:** Below average - Some relevant experience but significant gaps
- **0.61 - 0.75:** Good fit - Meets most requirements with minor gaps
- **0.76 - 0.85:** Strong fit - Exceeds requirements in several areas
- **0.86 - 1.00:** Excellent fit - Outstanding match with all criteria + bonus qualifications

### Project Score (1.0 - 5.0)

- **1.0 - 2.0:** Unacceptable - Major requirements not met
- **2.1 - 3.0:** Below expectations - Basic implementation with significant issues
- **3.1 - 3.7:** Meets expectations - Functional implementation with room for improvement
- **3.8 - 4.5:** Exceeds expectations - Well-executed with good practices
- **4.6 - 5.0:** Outstanding - Exceptional quality with creative enhancements

---

## Common Edge Cases to Handle

1. **Incomplete PDF Extraction:** Handle cases where PDF text extraction fails or returns garbled text
2. **LLM Timeout:** Implement retry with exponential backoff (max 3 attempts)
3. **Rate Limiting:** Queue jobs and process sequentially if needed
4. **Invalid JSON Response:** Parse with fallback regex extraction
5. **Missing Ground Truth:** Provide default context if vector DB query returns no results
6. **Concurrent Requests:** Ensure job queue handles multiple simultaneous evaluation requests
7. **Long Processing Times:** Maintain job status updates and prevent request timeouts

---

## Rubric Version

**Version:** 1.0  
**Last Updated:** 2025  
**Valid For:** Backend Developer / Product Engineer positions with AI/LLM components
