CREATE TABLE IF NOT EXISTS player_match_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    match_id UUID NOT NULL,
    hero_name VARCHAR(255) NOT NULL,
    kills INT NOT NULL,
    deaths INT NOT NULL,
    assists INT NOT NULL,
    player_rank_change INT NOT NULL,
    result VARCHAR(50) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_match
        FOREIGN KEY(match_id)
        REFERENCES matches(id)
        ON DELETE CASCADE
); 