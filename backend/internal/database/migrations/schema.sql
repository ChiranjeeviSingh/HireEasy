CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    "role" VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE jobs (
    id SERIAL PRIMARY KEY, --id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    job_title VARCHAR(255) NOT NULL, -- length validation in FE
    job_description TEXT NOT NULL,
    job_status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive
    skills_required VARCHAR[] NOT NULL, -- CHECK (array_length(skills_required, 1) > 0), can vaidate in FE
    attributes JSONB, --FE Q&A dump
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


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
);


-- Add indexes for common queries
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_id ON jobs(job_id);
CREATE INDEX idx_jobs_status ON jobs(job_status);
CREATE INDEX idx_jobs_title ON jobs(job_title);