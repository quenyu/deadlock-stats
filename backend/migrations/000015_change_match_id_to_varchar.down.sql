ALTER TABLE player_match_stats DROP CONSTRAINT fk_match;

ALTER TABLE player_match_stats ALTER COLUMN match_id TYPE UUID USING match_id::uuid;
ALTER TABLE matches ALTER COLUMN id TYPE UUID USING id::uuid;

ALTER TABLE matches ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE player_match_stats ADD CONSTRAINT fk_match FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE; 