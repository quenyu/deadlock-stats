CREATE TABLE crosshair_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    crosshair_id UUID NOT NULL REFERENCES crosshairs(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, crosshair_id)
);

CREATE INDEX idx_crosshair_likes_user_id ON crosshair_likes(user_id);
CREATE INDEX idx_crosshair_likes_crosshair_id ON crosshair_likes(crosshair_id);