# CV Scoring Rubric

## Overview
This document defines the standardized parameters and scoring methodology for evaluating candidate CVs against job requirements. Each parameter is scored on a scale of 1-5, with specific criteria for each level.

---

## Scoring Parameters

### 1. Technical Skills Match (Weight: 40%)

**Description**: Measures the alignment between the candidate's technical skills and the job requirements, focusing on backend technologies, databases, APIs, cloud platforms, and AI/LLM exposure.

**Scoring Guide**:

| Score | Criteria | Description |
|-------|----------|-------------|
| **5** | Excellent match + AI/LLM exposure | All required skills present (backend frameworks, databases, APIs, cloud). Strong AI/LLM integration experience. Demonstrates expertise in modern tech stack. |
| **4** | Strong match | Most required skills present with good depth. Some AI/LLM exposure. Solid backend and cloud experience. |
| **3** | Partial match | Covers core backend skills but missing some key areas. Limited or no AI/LLM experience. Moderate depth. |
| **2** | Few overlaps | Only basic backend skills present. Significant gaps in databases, cloud, or APIs. No AI/LLM experience. |
| **1** | Irrelevant skills | Skills do not align with job requirements. Missing critical technical competencies. |

**Evaluation Criteria**:
- **Backend Frameworks**: Experience with Express, FastAPI, Django, Rails, or similar
- **Databases**: SQL (PostgreSQL, MySQL) and NoSQL (MongoDB, Redis) proficiency
- **APIs**: RESTful API design and implementation
- **Cloud Platforms**: AWS, GCP, or Azure experience
- **AI/LLM**: OpenAI, Anthropic, Gemini APIs; prompt engineering; RAG systems
- **Vector Databases**: ChromaDB, Pinecone, Qdrant, Weaviate
- **Additional Skills**: Docker, Kubernetes, message queues, caching, monitoring

---

### 2. Experience Level (Weight: 25%)

**Description**: Evaluates the candidate's years of professional experience and the complexity of projects they have worked on.

**Scoring Guide**:

| Score | Criteria | Description |
|-------|----------|-------------|
| **5** | 5+ years / high-impact projects | 5+ years of backend development. Led or significantly contributed to large-scale, complex systems. Proven track record of shipping production applications. |
| **4** | 3-4 years solid track record | 3-4 years with consistent backend development experience. Multiple mid-to-large scale projects. Clear progression in responsibilities. |
| **3** | 2-3 years with mid-scale projects | 2-3 years of experience. Worked on moderately complex projects. Demonstrates growing technical maturity. |
| **2** | 1-2 years | 1-2 years of professional experience. Limited project complexity. Still developing core skills. |
| **1** | <1 year / trivial projects | Less than 1 year of experience or only worked on very simple projects. Limited exposure to real-world systems. |

**Evaluation Criteria**:
- Total years of professional backend development experience
- Complexity and scale of projects (users, data volume, traffic)
- Progression of responsibilities over time
- Variety of technologies and domains worked in
- Leadership or mentorship experience

---

### 3. Relevant Achievements (Weight: 20%)

**Description**: Assesses the measurable impact of the candidate's past work, including contributions to system performance, scalability, adoption, or business outcomes.

**Scoring Guide**:

| Score | Criteria | Description |
|-------|----------|-------------|
| **5** | Major measurable impact | Significant, quantifiable achievements (e.g., "reduced API latency by 80%", "scaled system to 1M+ users", "increased throughput by 10x"). Clear business impact. |
| **4** | Significant contributions | Multiple notable achievements with measurable outcomes. Demonstrates technical excellence and impact. |
| **3** | Some measurable outcomes | A few quantified achievements. Shows awareness of performance and quality improvements. |
| **2** | Minimal improvements | Limited or vague achievements. Few measurable outcomes or metrics provided. |
| **1** | No clear achievements | No specific achievements mentioned. Generic responsibilities without impact statements. |

**Evaluation Criteria**:
- **Performance Improvements**: Latency reduction, throughput increase, optimization results
- **Scalability**: System growth metrics (users, requests, data volume)
- **Reliability**: Uptime improvements, bug reduction, error rate decrease
- **Adoption**: Feature adoption rates, user growth, customer satisfaction
- **Cost Optimization**: Infrastructure cost savings, efficiency gains
- **Innovation**: Patents, open-source contributions, technical publications

---

### 4. Cultural / Collaboration Fit (Weight: 15%)

**Description**: Evaluates the candidate's communication skills, learning mindset, teamwork, and leadership capabilities as demonstrated in their CV.

**Scoring Guide**:

| Score | Criteria | Description |
|-------|----------|-------------|
| **5** | Excellent and well-demonstrated | Strong evidence of leadership, mentorship, cross-functional collaboration. Excellent communication (presentations, documentation, teaching). Active learning (certifications, courses, community involvement). |
| **4** | Good | Clear examples of teamwork and collaboration. Good communication skills. Demonstrates learning mindset through various activities. |
| **3** | Average | Some mention of team projects or collaboration. Basic communication and learning indicators present. |
| **2** | Minimal | Limited evidence of collaboration or communication skills. Few indicators of learning or growth mindset. |
| **1** | Not demonstrated | No evidence of teamwork, communication, or continuous learning. CV focused only on technical tasks. |

**Evaluation Criteria**:
- **Communication**: Technical writing, presentations, documentation, teaching/training
- **Collaboration**: Cross-functional projects, team achievements, pair programming
- **Learning Mindset**: Certifications, courses, workshops, conference attendance
- **Leadership**: Mentoring, technical leadership, project ownership
- **Community**: Open-source contributions, blog posts, speaking engagements
- **Soft Skills**: Problem-solving approach, adaptability, initiative

---

## Calculation Methodology

### Individual Parameter Score
Each parameter is scored from 1 to 5 based on the criteria above.

### Weighted Score Calculation
```
Weighted Score = (Technical Skills × 0.40) + (Experience Level × 0.25) + 
                 (Relevant Achievements × 0.20) + (Cultural Fit × 0.15)
```

### CV Match Rate Conversion
The weighted score (1-5 scale) is converted to a match rate (0-1 scale):
```
CV Match Rate = Weighted Score × 0.2
```

**Example**:
- Technical Skills: 4 (Strong match)
- Experience Level: 5 (5+ years)
- Relevant Achievements: 3 (Some measurable outcomes)
- Cultural Fit: 4 (Good)

```
Weighted Score = (4 × 0.40) + (5 × 0.25) + (3 × 0.20) + (4 × 0.15)
                = 1.6 + 1.25 + 0.6 + 0.6
                = 4.05

CV Match Rate = 4.05 × 0.2 = 0.81 (81% match)
```

---

## CV Feedback Guidelines

The AI evaluation should provide detailed feedback that includes:

1. **Strengths**: Highlight the candidate's strongest areas aligned with job requirements
2. **Gaps**: Identify missing skills or experience areas
3. **Specific Examples**: Reference concrete achievements or experiences from the CV
4. **Recommendations**: Suggest areas for improvement or clarification in interviews

**Example Feedback**:
> "Strong in backend development with 5 years of experience using Python/Django and Node.js/Express. Excellent cloud experience with AWS (EC2, Lambda, RDS). Demonstrated measurable impact with 60% API latency reduction and system scaling to 500K users. However, limited AI/LLM integration experience—only mentions basic OpenAI API usage. No evidence of RAG implementation or vector database work. Recommend exploring prompt engineering and RAG architectures to strengthen AI capabilities."

---

## Scoring Consistency Guidelines

To ensure consistent and reliable scoring:

1. **Evidence-Based**: Score based only on information explicitly stated in the CV
2. **No Assumptions**: Do not infer skills or experience not mentioned
3. **Comparative Analysis**: Compare candidate skills against job description requirements
4. **Objective Metrics**: Prioritize quantifiable achievements over subjective claims
5. **Balanced View**: Consider both depth (expertise level) and breadth (variety of skills)

---

## Edge Cases and Special Considerations

### Career Transitions
- If candidate is transitioning from related field (e.g., frontend to backend), evaluate transferable skills carefully
- Consider relevant projects or self-learning efforts

### Recent Graduates
- For candidates with <2 years experience, place more weight on project complexity and learning ability
- Consider internships, personal projects, and academic achievements

### Senior/Lead Roles
- For candidates with 7+ years, expect evidence of technical leadership and architectural decisions
- Look for mentorship, system design ownership, and broader impact

### AI/LLM Specialists
- If role heavily emphasizes AI/LLM work, consider increasing Technical Skills weight
- Evaluate specific AI competencies (prompt engineering, fine-tuning, RAG, agents)

---

## Version History
- **v1.0** - Initial scoring rubric for CV evaluation