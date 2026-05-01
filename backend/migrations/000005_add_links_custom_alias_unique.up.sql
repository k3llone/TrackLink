CREATE UNIQUE INDEX IF NOT EXISTS idx_links_custom_alias_unique
ON links (custom_alias)
WHERE custom_alias IS NOT NULL;
