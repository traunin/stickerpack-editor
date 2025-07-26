CREATE TABLE IF NOT EXISTS stickerpacks (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    is_public BOOLEAN NOT NULL
)
