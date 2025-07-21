ALTER TABLE player_match_stats
ADD CONSTRAINT unique_user_match UNIQUE (user_id, match_id); 