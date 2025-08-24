-- Insert initial loyalty tiers
INSERT INTO loyalty_tiers (id, name, min_points, description, benefits) VALUES
(1, 'Bronze', 0, 'Initial tier for all new users.', '5% discount on selected items'),
(2, 'Silver', 100, 'Achieved after accumulating 100 points.', '10% discount, early access to sales'),
(3, 'Gold', 500, 'Achieved after accumulating 500 points.', '15% discount, dedicated support, birthday gift'),
(4, 'Platinum', 1000, 'Achieved after accumulating 1000 points.', '20% discount, personal shopper, exclusive events');

-- Update existing users with initial loyalty status and points, and create user_loyalty entries
DO $$
DECLARE
    user_record RECORD;
    bronze_tier_id INT;
BEGIN
    -- Get the ID for the 'Bronze' tier
    SELECT id INTO bronze_tier_id FROM loyalty_tiers WHERE name = 'Bronze';

    FOR user_record IN SELECT id FROM users LOOP
        -- Update users table
        UPDATE users SET loyalty_status = 'Bronze', current_points = 0 WHERE id = user_record.id;

        -- Insert into user_loyalty table
        INSERT INTO user_loyalty (user_id, current_points, current_tier_id, last_activity_at)
        VALUES (user_record.id, 0, bronze_tier_id, CURRENT_TIMESTAMP);

        -- Add a welcome loyalty activity
        INSERT INTO loyalty_activities (user_id, type, description, created_at)
        VALUES (user_record.id, 'welcome_bonus', 'Received welcome bonus for joining loyalty program', CURRENT_TIMESTAMP);
    END LOOP;
END
$$;
