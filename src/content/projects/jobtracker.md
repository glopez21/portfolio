---
title: jobtracker
pillar: python
tagline: Job application intelligence — resume parsing, semantic match scoring, and
  AI cover letters.
status: dev
order: 4
featured: false
path: jobtracker
repo: https://git.4rch3.io/61-72-6b-68-e9/jobtracker
stack:
- Python
- FastAPI
- HTMX
- spaCy
- SQLite
highlights:
- Resume and job description parsing with spaCy NER
- Semantic match scoring against saved applications
- AI-generated cover letter drafts per JD
- HTMX-powered dashboard with real-time updates
---

jobtracker is a job application management tool with an intelligence
layer — not just a spreadsheet of applications, but a system that
analyzes each job description against your resume and generates
customized cover letter drafts.

The NLP pipeline uses spaCy for named entity recognition, extracting
skills, tools, and experience from both resumes and job descriptions.
The semantic match score quantifies how well your profile aligns with
each position, surfacing high-fit opportunities and highlighting gaps.

The HTMX frontend provides a live dashboard: saved jobs with match
scores, generated cover letters, and application status — all updated
in-place without full page reloads. The project is particularly
meta-valuable: it's a tool you use while job-hunting, built to
demonstrate the exact skills the jobs require.