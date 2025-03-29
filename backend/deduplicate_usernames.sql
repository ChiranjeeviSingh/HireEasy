-- Make usernames unique by adding a suffix to duplicates
DO $$
DECLARE
   r RECORD;
   counter INT;
BEGIN
   -- Loop through each duplicated username
   FOR r IN (
      SELECT username 
      FROM users 
      GROUP BY username 
      HAVING COUNT(*) > 1
   ) LOOP
      counter := 1;
      
      -- For each duplicate username, update all but the first instance with a numbered suffix
      UPDATE users
      SET username = username || '_' || id
      WHERE id IN (
         SELECT id
         FROM users
         WHERE username = r.username
         ORDER BY id
         OFFSET 1  -- Skip the first one
      );
   END LOOP;
END $$;

-- Verify results
SELECT username, COUNT(*)
FROM users
GROUP BY username
HAVING COUNT(*) > 1; 