DROP INDEX IF EXISTS idx_urls_short_url;
CREATE INDEX idx_urls_short_url ON urls(short_url);
