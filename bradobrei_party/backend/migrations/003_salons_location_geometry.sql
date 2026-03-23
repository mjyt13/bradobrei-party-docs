CREATE EXTENSION IF NOT EXISTS postgis;

ALTER TABLE salons
    ALTER COLUMN location TYPE geometry(POINT, 4326)
    USING CASE
        WHEN location IS NULL THEN NULL
        WHEN pg_typeof(location)::text = 'text' THEN ST_GeomFromText(location, 4326)
        ELSE location::geometry
    END;

DROP INDEX IF EXISTS idx_salons_location;
CREATE INDEX IF NOT EXISTS idx_salons_location ON salons USING GIST(location);
