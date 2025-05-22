CREATE TABLE
    IF NOT EXISTS post_privacy (
        post_id TEXT NOT NULL,
        chosen_id TEXT NOT NULL,
        PRIMARY KEY (post_id, chosen_id),
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        FOREIGN KEY (chosen_id) REFERENCES users (id) ON DELETE CASCADE
    );