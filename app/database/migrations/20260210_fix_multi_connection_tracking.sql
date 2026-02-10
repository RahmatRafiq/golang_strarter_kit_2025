-- +++ UP Migration
-- Fix multi-connection tracking by adding connection_name column
-- This ensures migrations and seeds are tracked independently per connection
--
-- ⚠️ IMPORTANT NOTE:
-- This migration file CANNOT be run through the normal migration system due to
-- a chicken-and-egg problem: the migration system queries 'connection_name' column
-- which doesn't exist yet. This migration was applied manually using direct SQL.
--
-- For fresh installations, the schema is already correct (see ensureMigrationsTable
-- and ensureSeedsTable in app/database/migration_manager.go and seeder_manager.go).
--
-- This file is kept for:
-- 1. Historical documentation of the schema change
-- 2. Reference for understanding the fix (see Issue #23)
-- 3. Manual application if needed on existing databases
--
-- To apply manually:
-- mysql -u root -p database_name < this_file.sql
--
-- Related: Issue #23, Issue #16

-- Step 1: Add connection_name column to migrations table
-- Check if column exists first (for idempotency)
SET @db_type = (SELECT IF(COUNT(*) > 0, 'mysql', 'postgres') 
                FROM information_schema.tables 
                WHERE table_schema = DATABASE() 
                LIMIT 1);

-- For MySQL/MariaDB
ALTER TABLE migrations 
ADD COLUMN IF NOT EXISTS connection_name VARCHAR(50) NOT NULL DEFAULT 'mysql' AFTER id;

-- Update existing records to use 'mysql' as connection name
UPDATE migrations SET connection_name = 'mysql' WHERE connection_name = '';

-- Add unique constraint to prevent duplicate migrations per connection
-- First drop the old constraint if exists
ALTER TABLE migrations DROP INDEX IF EXISTS idx_migrations_filename;

-- Add new unique constraint with connection_name
ALTER TABLE migrations 
ADD UNIQUE KEY unique_migration (connection_name, filename);

-- Step 2: Add connection_name column to seeds table
ALTER TABLE seeds 
ADD COLUMN IF NOT EXISTS connection_name VARCHAR(50) NOT NULL DEFAULT 'mysql' AFTER id;

-- Update existing seed records
UPDATE seeds SET connection_name = 'mysql' WHERE connection_name = '';

-- Add unique constraint for seeds
ALTER TABLE seeds DROP INDEX IF EXISTS idx_seeds_filename;

ALTER TABLE seeds 
ADD UNIQUE KEY unique_seed (connection_name, filename);

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_migrations_connection ON migrations(connection_name);
CREATE INDEX IF NOT EXISTS idx_seeds_connection ON seeds(connection_name);

-- --- DOWN Migration
-- Rollback: Remove connection_name column and constraints

-- Remove unique constraints
ALTER TABLE migrations DROP INDEX IF EXISTS unique_migration;
ALTER TABLE seeds DROP INDEX IF EXISTS unique_seed;

-- Remove performance indexes
ALTER TABLE migrations DROP INDEX IF EXISTS idx_migrations_connection;
ALTER TABLE seeds DROP INDEX IF EXISTS idx_seeds_connection;

-- Remove connection_name columns
ALTER TABLE migrations DROP COLUMN IF EXISTS connection_name;
ALTER TABLE seeds DROP COLUMN IF EXISTS connection_name;

-- Restore old unique constraints (optional, if they existed)
-- ALTER TABLE migrations ADD UNIQUE KEY idx_migrations_filename (filename);
-- ALTER TABLE seeds ADD UNIQUE KEY idx_seeds_filename (filename);
