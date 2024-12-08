CREATE TABLE order_statuses
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE orders
(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    status INTEGER NOT NULL REFERENCES order_statuses(id) DEFAULT 1,
    total DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX orders_user_id_index ON orders(user_id);
CREATE INDEX orders_id_index ON orders(id);

CREATE TABLE order_items
(
    order_id UUID NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (order_id, product_id)
);

INSERT INTO order_statuses (name) VALUES 
('PENDING'),
('APPROVED'),
('REJECTED');