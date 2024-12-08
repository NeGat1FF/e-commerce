CREATE TABLE payment_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE payments (
    id UUID PRIMARY KEY,
    stripe_id VARCHAR(255) NOT NULL,
    order_id UUID NOT NULL,
    status INT NOT NULL REFERENCES payment_statuses(id),
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO payment_statuses (name) VALUES 
('PENDING'), 
('CONFIRMED'), 
('FAILED');

