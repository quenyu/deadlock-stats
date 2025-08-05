CREATE TABLE crosshairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    settings JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    view_count INT NOT NULL DEFAULT 0,
    likes_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_crosshairs_author_id ON crosshairs(author_id);
CREATE INDEX idx_crosshairs_created_at ON crosshairs(created_at DESC);
CREATE INDEX idx_crosshairs_likes_count ON crosshairs(likes_count DESC);