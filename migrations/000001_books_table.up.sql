CREATE TABLE IF NOT EXISTS books (
                                     id SERIAL PRIMARY KEY,
                                     title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publish_date TIMESTAMP NOT NULL DEFAULT NOW(),
    rating INTEGER CHECK (rating >= 0 AND rating <= 5) DEFAULT 0
    );

CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_books_title ON books(title);