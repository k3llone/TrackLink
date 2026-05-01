DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'links_owner_id_fkey'
    ) THEN
        ALTER TABLE links
        ADD CONSTRAINT links_owner_id_fkey
        FOREIGN KEY (owner_id) REFERENCES users(id)
        ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_links_owner_id ON links (owner_id);
