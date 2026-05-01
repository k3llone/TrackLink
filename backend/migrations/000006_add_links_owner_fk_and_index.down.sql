DROP INDEX IF EXISTS idx_links_owner_id;

ALTER TABLE links
DROP CONSTRAINT IF EXISTS links_owner_id_fkey;
