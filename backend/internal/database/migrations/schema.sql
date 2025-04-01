CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- DROP TABLE IF EXISTS application_form CASCADE;
-- DROP TABLE IF EXISTS form_templates CASCADE;
-- DROP TABLE IF EXISTS jobs CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS profiles CASCADE;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    "role" VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
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
    id SERIAL PRIMARY KEY,
    form_template_id VARCHAR(255) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    fields JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS application_form (
    form_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),   -- Auto-generating unique UUID
    job_id INT NOT NULL,                                     -- id of job table
    form_id INT NOT NULL,                                    -- id of form_template table
    status VARCHAR(50) NOT NULL DEFAULT 'active',            -- active, inactive
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,      -- Date when the form is created
    FOREIGN KEY (job_id) REFERENCES jobs(id),                -- Foreign key reference to the jobs table
    FOREIGN KEY (form_id) REFERENCES form_templates(id)      -- Foreign key reference to the form_templates table
);

CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_title VARCHAR(255),
    years_of_experience INT,
    areas_of_expertise TEXT[],
    phone_number VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

---remove job_submissions after main is pulled
CREATE TABLE IF NOT EXISTS job_submissions (
    id SERIAL PRIMARY KEY,
    form_uuid VARCHAR(255) NOT NULL REFERENCES form_templates(form_template_id) ON DELETE CASCADE,
    job_id VARCHAR(255) NOT NULL REFERENCES jobs(job_id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    form_data JSONB NOT NULL, -- Stores user responses dynamically
    resume_url TEXT NOT NULL, -- Store S3 URL instead of local path
    ats_score INTEGER NOT NULL DEFAULT 0, -- ATS ranking score (0-100)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS availabilities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "date" DATE NOT NULL, -- Req Expected format: YYYY-MM-DD
    from_time TIME(0) NOT NULL,
    to_time TIME(0) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,
    hr_user_id INT REFERENCES users(id) ON DELETE SET NULL,
    job_submission_id INT NOT NULL REFERENCES job_submissions(id),
    interviewer_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    availability_id INT NOT NULL REFERENCES availabilities(id) ON DELETE CASCADE,
    feedback TEXT DEFAULT NULL,
    "status" VARCHAR(50) DEFAULT 'scheduled', --or completed 
    UNIQUE (availability_id) -- Ensures an availability slot can only be booked once(this is ensured from application also)
);


-- Add indexes for common queries
CREATE INDEX IF NOT EXISTS idx_jobs_user_id ON jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_jobs_id ON jobs(job_id);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(job_status);
CREATE INDEX IF NOT EXISTS idx_jobs_title ON jobs(job_title);

CREATE INDEX IF NOT EXISTS idx_availabilities_user ON availabilities (user_id, date);
CREATE INDEX IF NOT EXISTS idx_availabilities_time ON availabilities (date, from_time, to_time);
CREATE INDEX IF NOT EXISTS idx_interview_availabilities ON interviews (availability_id);
