CREATE TABLE IF NOT EXISTS content_tags (
    tag_id INT NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,

    PRIMARY KEY (tag_id, content_type, content_id),
    CONSTRAINT fk_tag
        FOREIGN KEY(tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
    -- We don't add a foreign key to content_id because it can reference multiple tables (builds, crosshairs).
    -- This logic will be handled at the application level.
);

CREATE INDEX IF NOT EXISTS idx_content_tags_content ON content_tags(content_type, content_id);
