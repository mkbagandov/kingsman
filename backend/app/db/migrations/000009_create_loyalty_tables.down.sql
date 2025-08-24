-- Remove loyalty_status and current_points from users table
ALTER TABLE users
DROP COLUMN loyalty_status,
DROP COLUMN current_points;

-- Drop user_loyalty table
DROP TABLE IF EXISTS user_loyalty;

-- Drop loyalty_activities table
DROP TABLE IF EXISTS loyalty_activities;

-- Drop loyalty_points table
DROP TABLE IF EXISTS loyalty_points;

-- Drop loyalty_tiers table
DROP TABLE IF EXISTS loyalty_tiers;
