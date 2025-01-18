-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS authors (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "name" VARCHAR(255),
  "created_at" TIMESTAMPTZ DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authors;
-- +goose StatementEnd