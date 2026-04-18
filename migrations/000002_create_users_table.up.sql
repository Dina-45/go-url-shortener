CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     created_at TIMESTAMP,
                                     updated_at TIMESTAMP,
                                     deleted_at TIMESTAMP,
                                     username VARCHAR(100) UNIQUE NOT NULL,
                                     password TEXT NOT NULL
);