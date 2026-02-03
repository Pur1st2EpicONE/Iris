CREATE TABLE IF NOT EXISTS links (
    id             INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    original_link  TEXT NOT NULL,
    short_link     TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS visits (
    id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    link_id     INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    user_agent  TEXT,
    visited_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_visits_link_id ON visits(link_id);
CREATE INDEX IF NOT EXISTS idx_visits_link_id_visited_at_user_agent ON visits(link_id, visited_at, user_agent);
