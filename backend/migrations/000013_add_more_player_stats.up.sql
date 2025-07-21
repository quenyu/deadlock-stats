ALTER TABLE player_stats
ADD COLUMN max_kills_in_match INT NOT NULL DEFAULT 0,
ADD COLUMN avg_damage_per_match REAL NOT NULL DEFAULT 0,
ADD COLUMN avg_objectives_per_match REAL NOT NULL DEFAULT 0,
ADD COLUMN avg_souls_per_min REAL NOT NULL DEFAULT 0; 