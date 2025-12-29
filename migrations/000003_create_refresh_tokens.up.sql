CREATE TABLE refresh_tokens (
                                id SERIAL PRIMARY KEY,
                                user_id INT REFERENCES users(id) ON DELETE CASCADE,
                                token VARCHAR(255) NOT NULL,
                                expires_at TIMESTAMP NOT NULL
);