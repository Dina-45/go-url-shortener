CREATE TABLE IF NOT EXISTS urls (
                                    id SERIAL PRIMARY KEY,
                                    created_at TIMESTAMP,
                                    updated_at TIMESTAMP,
                                    deleted_at TIMESTAMP,
                                    original_url TEXT NOT NULL,
                                    short_code VARCHAR(50) UNIQUE NOT NULL
);