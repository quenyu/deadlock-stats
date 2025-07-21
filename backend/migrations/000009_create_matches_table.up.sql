CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    map_name VARCHAR(255) NOT NULL,
    duration_minutes INT NOT NULL,
    match_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 