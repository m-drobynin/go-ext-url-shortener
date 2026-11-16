DROP INDEX IF EXISTS idx_urls_short_url;
CREATE UNIQUE INDEX idx_urls_short_url ON urls(short_url);