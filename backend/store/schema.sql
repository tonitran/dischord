-- DisChord schema. Applied automatically on startup via store.ApplySchema.

CREATE TABLE IF NOT EXISTS users (
    id         TEXT PRIMARY KEY,
    username   TEXT NOT NULL DEFAULT '',
    email      TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS servers (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL DEFAULT '',
    owner_id   TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id         TEXT PRIMARY KEY,
    server_id  TEXT NOT NULL DEFAULT '',
    author_id  TEXT NOT NULL DEFAULT '',
    title      TEXT NOT NULL DEFAULT '',
    body       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS votes (
    post_id   TEXT NOT NULL,
    author_id TEXT NOT NULL,
    vote      INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (post_id, author_id)
);

CREATE TABLE IF NOT EXISTS friends (
    user_id   TEXT NOT NULL,
    friend_id TEXT NOT NULL,
    PRIMARY KEY (user_id, friend_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id         TEXT PRIMARY KEY,
    server_id  TEXT NOT NULL DEFAULT '',
    author_id  TEXT NOT NULL DEFAULT '',
    content    TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
