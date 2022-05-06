DROP TABLE IF EXISTS urls;

CREATE TABLE urls (
    fullurl TEXT NOT NULL UNIQUE,
    tinyurl TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS index_fullurl ON urls USING HASH (fullurl);
CREATE INDEX IF NOT EXISTS index_tinyurl ON urls USING HASH (tinyurl);