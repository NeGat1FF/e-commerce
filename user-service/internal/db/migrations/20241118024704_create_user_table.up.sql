CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone VARCHAR(15) UNIQUE,
    password VARCHAR(255) NOT NULL,
    refresh_tokens VARCHAR(256)[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_email_idx ON users(email);

CREATE INDEX users_token_idx ON users(refresh_tokens);

CREATE TABLE email_verifications (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX email_verifications_token_idx ON email_verifications(token);

CREATE TABLE password_resets (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(128) NOT NULL UNIQUE,       
    expires_at TIMESTAMP NOT NULL    
);

CREATE INDEX password_resets_token_idx ON password_resets(token);