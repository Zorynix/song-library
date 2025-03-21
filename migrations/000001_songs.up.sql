CREATE SCHEMA IF NOT EXISTS library;

CREATE TABLE library.songs (
    id SERIAL PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    release_date VARCHAR(50) NOT NULL,
    text TEXT NOT NULL,
    link VARCHAR(255) NOT NULL
);