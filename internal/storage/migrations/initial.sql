CREATE TABLE IF NOT EXISTS "urls" (
  "short_url" varchar PRIMARY KEY,
  "original_url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "user_id" varchar,
  "is_deleted" int DEFAULT (0)
);

CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON urls (original_url);