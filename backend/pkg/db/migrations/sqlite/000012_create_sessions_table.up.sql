CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,                          
    user_id INTEGER NOT NULL,                     
    expires_at TEXT NOT NULL,                     
    created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);