-- Haoma Database Initialization Script
-- This script sets up the initial database structure

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create indexes for better performance
-- (GORM will auto-create tables, but we can optimize with indexes)

-- Note: This file runs automatically when PostgreSQL starts for the first time
-- Additional indexes will be added here as needed

-- Example: After tables are created by GORM, you might want to add:
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sessions_player_id ON sessions(player_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_sessions_started_at ON sessions(started_at);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_leaderboard_entries_score ON leaderboard_entries(final_score DESC, completion_time ASC);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_attempts_session_id ON attempts(session_id);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_questions_category_id ON questions(category_id);
