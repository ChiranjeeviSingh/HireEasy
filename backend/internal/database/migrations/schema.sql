CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY, --id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id VARCHAR(255) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    job_title VARCHAR(255) NOT NULL, -- length validation in FE
    job_description TEXT NOT NULL,
    job_status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive
    skills_required VARCHAR[] NOT NULL, -- CHECK (array_length(skills_required, 1) > 0), can vaidate in FE
    attributes JSONB, --FE Q&A dump
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


CREATE TABLE IF NOT EXISTS form_templates (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    job_id VARCHAR(255) NOT NULL REFERENCES jobs(job_id) ON DELETE CASCADE,
    fields JSONB NOT NULL, -- Dynamic form fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_submissions (
    id SERIAL PRIMARY KEY,
    form_uuid UUID NOT NULL REFERENCES form_templates(id) ON DELETE CASCADE,
    job_id VARCHAR(255) NOT NULL REFERENCES jobs(job_id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    form_data JSONB NOT NULL, -- Stores user responses dynamically
    resume_url TEXT NOT NULL, -- Store S3 URL instead of local path
    ats_score INTEGER NOT NULL DEFAULT 0, -- ATS ranking score (0-100)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Add indexes for common queries

CREATE INDEX idx_job_submissions_job ON job_submissions(form_uuid);


-- Create index only if it does not exist
DO $$ 
BEGIN 
  -- Ensure user ID index exists in jobs table
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_jobs_user_id') THEN
    CREATE INDEX idx_jobs_user_id ON jobs(user_id);
  END IF;

  -- Ensure job submissions index exists
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_job_submissions_job') THEN
    CREATE INDEX idx_job_submissions_job ON job_submissions(job_id);
  END IF;

  -- Fix incorrect index name (was idx_jobs_submissions)
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_jobs_id') THEN
    CREATE INDEX idx_jobs_id ON jobs(job_id);
  END IF;

  -- Ensure job status index exists
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_jobs_status') THEN
    CREATE INDEX idx_jobs_status ON jobs(job_status);
  END IF;

  -- Ensure job title index exists
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_jobs_title') THEN
    CREATE INDEX idx_jobs_title ON jobs(job_title);
  END IF;

  -- Ensure ATS score index exists for job submissions
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_job_submissions_ats') THEN
    CREATE INDEX idx_job_submissions_ats ON job_submissions(ats_score DESC);
  END IF;
END $$;

-- Add triggers to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for jobs table
CREATE TRIGGER update_jobs_updated_at
BEFORE UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create trigger for users table
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create a unique constraint to prevent duplicate job submissions
ALTER TABLE job_submissions ADD CONSTRAINT unique_job_submission UNIQUE (job_id, email);
