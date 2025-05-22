CREATE TABLE IF NOT EXISTS posts (
    id Integer PRIMARY KEY autoincrement,
    user_id TEXT NOT NULL,
    group_id TEXT,
    content TEXT NOT NULL,
    image TEXT,
    privacy TEXT CHECK(privacy IN ('public', 'almost_private', 'private', '')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
