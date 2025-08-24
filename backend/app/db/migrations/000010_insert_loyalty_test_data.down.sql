-- Remove test loyalty data

-- Delete loyalty activities
DELETE FROM loyalty_activities WHERE type = 'welcome_bonus';

-- Delete user_loyalty entries
DELETE FROM user_loyalty WHERE current_points = 0;

-- Reset loyalty_status and current_points for all users
UPDATE users SET loyalty_status = '', current_points = 0;

-- Delete loyalty tiers
DELETE FROM loyalty_tiers WHERE id IN (1, 2, 3, 4);
