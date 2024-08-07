CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    user_name TEXT NOT NULL UNIQUE,
    login TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_user_name ON users(user_name);