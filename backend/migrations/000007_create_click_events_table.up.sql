CREATE TABLE IF NOT EXISTS click_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id UUID NOT NULL,
    clicked_at TIMESTAMPTZ NOT NULL,
    referrer TEXT NULL,
    user_agent TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT click_events_link_id_fkey
        FOREIGN KEY (link_id) REFERENCES links(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_click_events_link_id ON click_events (link_id);
CREATE INDEX IF NOT EXISTS idx_click_events_clicked_at ON click_events (clicked_at DESC);
