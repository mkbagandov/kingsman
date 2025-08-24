-- Add loyalty_status and current_points to users table
ALTER TABLE users
ADD COLUMN loyalty_status VARCHAR(50) DEFAULT 'Bronze' NOT NULL,
ADD COLUMN current_points INT DEFAULT 0 NOT NULL;

-- Create loyalty_tiers table
CREATE TABLE loyalty_tiers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    min_points INT NOT NULL,
    description TEXT,
    benefits TEXT
);

-- Create loyalty_points table
CREATE TABLE loyalty_points (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    points INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loyalty_activities table
CREATE TABLE loyalty_activities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create user_loyalty table to aggregate loyalty information
CREATE TABLE user_loyalty (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_points INT DEFAULT 0 NOT NULL,
    current_tier_id INT REFERENCES loyalty_tiers(id) ON DELETE SET NULL,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
