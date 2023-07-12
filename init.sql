CREATE DATABASE my_database;
CREATE USER my_user WITH ENCRYPTED PASSWORD 'my_password';
GRANT ALL PRIVILEGES ON DATABASE my_database TO my_user;

DROP TABLE IF EXISTS songs;

-- Crear la tabla songs
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    artist TEXT NOT NULL,
    duration TEXT,
    album TEXT,
    artwork TEXT,
    price TEXT,
    origin TEXT NOT NULL
);