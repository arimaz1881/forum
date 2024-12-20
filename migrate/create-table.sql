CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT (CURRENT_TIMESTAMP),
    oauth_provider TEXT,
    oauth_id TEXT,
    login TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    hashed_password TEXT,
    role TEXT NOT NULL,
    moderator_role_request BOOLEAN DEFAULT 'f'
);

INSERT INTO
    users (
        login,
        email,
        hashed_password,
        role
    )
VALUES
    (
        'admin',
        'admin@main.com',
        'todo:write_hash',
        'admin'
    );

CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT (CURRENT_TIMESTAMP),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    file_key TEXT,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE categories (id INTEGER PRIMARY KEY, title TEXT);

INSERT INTO
    categories (id, title)
VALUES
    (1, 'A'),
    (2, 'F'),
    (3, 'N'),
    (4, 'G'),
    (5, 'O');

CREATE TABLE posts_categories (
    post_id INTEGER NOT NULL,
    categoria_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, categoria_id),
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (categoria_id) REFERENCES categories (id)
);

CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT (CURRENT_TIMESTAMP),
    content TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE post_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER NOT NULL,
    action TEXT,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE (post_id, user_id)
);

CREATE TABLE comment_reactions (
    comment_id INTEGER,
    user_id INTEGER NOT NULL,
    action TEXT,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (comment_id) REFERENCES comments (id),
    UNIQUE (comment_id, user_id)
);

CREATE TABLE sessions (user_id INTEGER, token TEXT, expires_at TIME);

CREATE TABLE notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT (CURRENT_TIMESTAMP),
    post_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    action_id INTEGER UNIQUE,
    comment_id INTEGER UNIQUE,
    action_type TEXT NOT NULL,
    seen BOOLEAN DEFAULT 0 NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (action_id) REFERENCES post_reactions (id),
    FOREIGN KEY (comment_id) REFERENCES comments (id)
);