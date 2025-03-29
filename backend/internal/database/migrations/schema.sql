CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

<<<<<<< HEAD
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY, --id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id VARCHAR(255) NOT NULL UNIQUE,
=======
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY, --id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id VARCHAR(255) NOT NULL,
>>>>>>> 26faa985bf52fc2df5a97b64d9e77ee62aad23a0
    user_id INTEGER NOT NULL REFERENCES users(id),
    job_title VARCHAR(255) NOT NULL, -- length validation in FE
    job_description TEXT NOT NULL,
    job_status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive
    skills_required VARCHAR[] NOT NULL, -- CHECK (array_length(skills_required, 1) > 0), can vaidate in FE
    attributes JSONB, --FE Q&A dump
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


<<<<<<< HEAD
CREATE TABLE IF NOT EXISTS form_templates (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(job_id) ON DELETE CASCADE,
    fields JSONB NOT NULL, -- Dynamic form fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_submissions (
    id SERIAL PRIMARY KEY,
    form_uuid UUID NOT NULL REFERENCES form_templates(id) ON DELETE CASCADE,
    job_id UUID NOT NULL REFERENCES jobs(job_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    form_data JSONB NOT NULL, -- Stores user responses dynamically
    resume_url TEXT NOT NULL, -- Store S3 URL instead of local path
    ats_score INTEGER NOT NULL DEFAULT 0, -- ATS ranking score (0-100)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
=======
CREATE TABLE form_templates (
    id SERIAL PRIMARY KEY,
    form_template_id VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    fields JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE application_form (
    form_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),   -- Auto-generating unique UUID
    job_id INT NOT NULL,                                     -- id of job table
    form_id INT NOT NULL,                                    -- id of form_template table
    status VARCHAR(50) NOT NULL DEFAULT 'active',            -- active, inactive
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,      -- Date when the form is created
    FOREIGN KEY (job_id) REFERENCES jobs(id),                -- Foreign key reference to the jobs table
    FOREIGN KEY (form_id) REFERENCES form_templates(id)      -- Foreign key reference to the form_templates table
>>>>>>> 26faa985bf52fc2df5a97b64d9e77ee62aad23a0
);


-- Add indexes for common queries
<<<<<<< HEAD

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

=======
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_id ON jobs(job_id);
CREATE INDEX idx_jobs_status ON jobs(job_status);
CREATE INDEX idx_jobs_title ON jobs(job_title);
>>>>>>> 26faa985bf52fc2df5a97b64d9e77ee62aad23a0
