INSERT INTO users (id, phone_number, social_id, discount_level, progress_to_next_level, qr_code) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '1234567890', NULL, 0, 0.0, 'qr123'),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '0987654321', 'social_id_1', 1, 25.5, 'qr456');

INSERT INTO stores (id, name, address, location, phone) VALUES
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Kingsman Store 1', '123 Main St', '1.23,4.56', '555-1111'),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Kingsman Store 2', '456 Oak Ave', '7.89,10.11', '555-2222');

INSERT INTO notifications (id, user_id, type, title, message, created_at) VALUES
    ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'promotion', 'Flash Sale!', 'Get 20% off all items this weekend!', '2025-08-24T10:00:00Z'),
    ('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'new_arrival', 'New Collection!', 'Check out our latest arrivals in store.', '2025-08-24T11:30:00Z');
