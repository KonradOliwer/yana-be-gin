CREATE TABLE IF NOT EXISTS "notes"
(
    id      UUID PRIMARY KEY NOT NULL,
    name    VARCHAR(50) UNIQUE,
    content TEXT
);
