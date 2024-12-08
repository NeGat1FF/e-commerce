CREATE TABLE cart (
    user_id UUID,
    item_id INTEGER,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (user_id, item_id)
);