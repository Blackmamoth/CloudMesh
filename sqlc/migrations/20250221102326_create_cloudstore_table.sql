-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cloud_store(
  id UUID DEFAULT gen_random_uuid(),
  account_id UUID NOT NULL,
  provider_id TEXT NOT NULL,
  provider_file_id TEXT NOT NULL,
  file_name TEXT NOT NULL,
  file_mime_type TEXT NOT NULL,
  file_size INT NOT NULL,
  file_created_time TEXT,
  file_modified_time TEXT,
  file_thumbnail_link TEXT,
  file_extension TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  UNIQUE (provider_file_id),
  FOREIGN KEY (account_id) REFERENCES account(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cloud_store;
-- +goose StatementEnd
