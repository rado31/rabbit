CREATE TABLE IF NOT EXISTS clients (
    id         SERIAL PRIMARY KEY,
    surname    VARCHAR(100) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    age        INT          NOT NULL CHECK (age > 0 AND age <= 150),
    email      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
