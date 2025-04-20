-- Check if form_uuid column exists in job_submissions table
DO $$
BEGIN
    -- Check if the column exists
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'job_submissions' AND column_name = 'form_uuid'
    ) THEN
        -- Add the column as TEXT to be compatible with any ID type
        ALTER TABLE job_submissions ADD COLUMN form_uuid TEXT DEFAULT '00000000-0000-0000-0000-000000000000';
        
        -- Make it NOT NULL after adding with default
        ALTER TABLE job_submissions ALTER COLUMN form_uuid SET NOT NULL;
    END IF;
END $$; 