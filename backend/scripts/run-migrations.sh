#!/bin/sh
set -eu

: "${DATABASE_URL:?DATABASE_URL is required}"

psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
SQL

for migration in /migrations/*.up.sql; do
    version="$(basename "$migration" .up.sql)"
    applied="$(psql "$DATABASE_URL" -tAc "SELECT 1 FROM schema_migrations WHERE version = '$version'")"

    if [ "$applied" = "1" ]; then
        echo "Skipping migration $version"
        continue
    fi

    echo "Applying migration $version"
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -v migration_version="$version" <<SQL
BEGIN;
\\i $migration
INSERT INTO schema_migrations (version) VALUES (:'migration_version');
COMMIT;
SQL
done
