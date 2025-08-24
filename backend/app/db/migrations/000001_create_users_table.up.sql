CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    social_id VARCHAR(255),
    discount_level INTEGER NOT NULL DEFAULT 0,
    progress_to_next_level NUMERIC(5, 2) NOT NULL DEFAULT 0.0,
    qr_code VARCHAR(36) UNIQUE
);
