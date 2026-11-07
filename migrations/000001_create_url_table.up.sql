CREATE TABLE urls (
    uuid UUID PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(2550) NOT NULL
);

-- Базовый индекс для поиска по названию
CREATE INDEX idx_urls_short_url ON urls(short_url);
