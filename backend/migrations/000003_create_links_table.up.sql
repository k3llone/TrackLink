CREATE TABLE IF NOT EXISTS links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    code TEXT NOT NULL,
    custom_alias TEXT NULL,
    target_url TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    total_clicks BIGINT NOT NULL DEFAULT 0,
    last_clicked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);
