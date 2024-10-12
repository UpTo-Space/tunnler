CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE,
    email VARCHAR(100) UNIQUE,
    activated BOOLEAN NOT NULL DEFAULT FALSE,
    activation_code BIGINT,
    password_hash VARCHAR(100)
);