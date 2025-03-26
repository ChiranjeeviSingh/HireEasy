-- Migration to add unique constraint to username field in users table
-- This will be applied to existing tables

-- First, check if any duplicate usernames exist
DO $$
DECLARE
    duplicate_count INTEGER;
BEGIN
    -- Count duplicate usernames
    SELECT COUNT(*) INTO duplicate_count
    FROM (
        SELECT username, COUNT(*) 
        FROM users 
        GROUP BY username 
        HAVING COUNT(*) > 1
    ) AS duplicates;
    
    -- If duplicates exist, raise an error (this should be handled in application code)
    IF duplicate_count > 0 THEN
        RAISE NOTICE 'There are % duplicate usernames in the users table. Please resolve these before adding the constraint.', duplicate_count;
    END IF;
END $$;

-- Add the unique constraint (will fail if duplicates exist)
DO $$
BEGIN
    -- Check if constraint already exists
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'users_username_key' AND conrelid = 'users'::regclass
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
        RAISE NOTICE 'Unique constraint added to username field in users table.';
    ELSE
        RAISE NOTICE 'Unique constraint already exists on username field in users table.';
    END IF;
EXCEPTION
    WHEN duplicate_key_value_violates_unique_constraint THEN
        RAISE EXCEPTION 'Could not add unique constraint due to duplicate usernames.';
END $$; 