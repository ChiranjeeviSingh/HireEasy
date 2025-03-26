# Database Migrations

This directory contains SQL migration files for the HireEasy application.

## How to Apply Migrations

### Initial Schema

The `schema.sql` file contains the initial database schema. Run this first when setting up a new database:

```bash
psql -U yourusername -d app_db -f schema.sql
```

### Migration: Add Username Unique Constraint

The `add_username_unique_constraint.sql` file adds a unique constraint to the username field in the users table:

```bash
psql -U yourusername -d app_db -f add_username_unique_constraint.sql
```

**Note:** This migration will fail if there are duplicate usernames in the database. The script will output a notice if duplicates exist. You'll need to resolve these before the constraint can be applied.

## Migration Order

Apply migrations in the following order:

1. `schema.sql` (on initial setup only)
2. `add_username_unique_constraint.sql`
3. (additional migrations will be added here)

## Troubleshooting

If you encounter errors with uniqueness constraints, you can identify duplicate values with:

```sql
SELECT username, COUNT(*) 
FROM users 
GROUP BY username 
HAVING COUNT(*) > 1;
```

Then update or remove the duplicate entries before applying the migration. 